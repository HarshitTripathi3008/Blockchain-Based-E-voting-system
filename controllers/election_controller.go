package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"MAJOR-PROJECT/bindings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

var (
	l2ClientMu            sync.Mutex
	l2ClientShared        *ethclient.Client
	verifiedContractAddrs sync.Map
)

// Candidate is the representation returned to the client
type Candidate struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ImageHash    string   `json:"imageHash"`
	VoteCount    *big.Int `json:"voteCount"`
	Email        string   `json:"email"`
	ManifestoUrl string   `json:"manifestoUrl,omitempty"`
}

type BlockchainResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// writeJSONHeader sets common JSON + CORS headers
func writeJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// respond helpers for consistent responses
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	writeJSONHeader(w)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, BlockchainResponse{Status: "error", Message: message})
}

func GetVoterStatus(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	address := r.URL.Query().Get("address")
	email := r.URL.Query().Get("email")

	if address == "" || email == "" {
		respondError(w, http.StatusBadRequest, "Missing address or email")
		return
	}

	// 1. Check job queue first (most recent)
	job, _ := getVoteJobByVoter(address, email)
	if job != nil {
		respondJSON(w, http.StatusOK, BlockchainResponse{
			Status:  "success",
			Message: "Voter status found in job queue",
			Data: map[string]interface{}{
				"hasVoted": true,
				"status":   job.Status,
				"txHash":   job.TxHash,
			},
		})
		return
	}

	// 2. Check Voter registration status in DB (permanent)
	var voter Voter
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addrRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(address) + "$", "$options": "i"}
	emailRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(email) + "$", "$options": "i"}

	// Find voter who has a registration for this address
	err := voterCollection.FindOne(ctx, bson.M{
		"email": emailRegex,
		"registrations.election_address": addrRegex,
	}).Decode(&voter)

	if err == nil {
		for _, reg := range voter.Registrations {
			if strings.EqualFold(reg.ElectionAddress, address) &&
				(strings.EqualFold(reg.Status, "Voted") || strings.EqualFold(reg.Status, "voted")) {
				respondJSON(w, http.StatusOK, BlockchainResponse{
					Status:  "success",
					Message: "Voter has already voted",
					Data: map[string]interface{}{
						"hasVoted": true,
						"status":   "mined",
					},
				})
				return
			}
		}
	}

	respondJSON(w, http.StatusOK, BlockchainResponse{
		Status:  "success",
		Message: "Voter has not voted yet",
		Data: map[string]interface{}{
			"hasVoted": false,
		},
	})
}

// getClient connects to L2_NODE_URL (with timeout)
func getClient() (*ethclient.Client, error) {
	l2ClientMu.Lock()
	defer l2ClientMu.Unlock()

	if l2ClientShared != nil {
		return l2ClientShared, nil
	}

	nodeURL := strings.TrimSpace(os.Getenv("L2_NODE_URL"))
	nodeURL = strings.Trim(nodeURL, `"'`)
	if nodeURL == "" {
		log.Println("DEBUG: raw L2_NODE_URL from env = ''")
		return nil, fmt.Errorf("L2_NODE_URL not configured")
	}

	// If the value is a host:port without scheme, add http://
	if !strings.HasPrefix(nodeURL, "http://") && !strings.HasPrefix(nodeURL, "https://") && !strings.HasPrefix(nodeURL, "ws://") && !strings.HasPrefix(nodeURL, "wss://") {
		nodeURL = "http://" + nodeURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(ctx, nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum node %q: %w", nodeURL, err)
	}
	l2ClientShared = client
	return l2ClientShared, nil
}

// getAuth creates a transact opts using EVM_PRIVATE_KEY and L2_CHAIN_ID.
// It optionally uses the shared client to fetch a suggested gas price once per request.
func getAuth(client *ethclient.Client) (*bind.TransactOpts, error) {
	// 1. Get Private Key
	priv := strings.TrimSpace(os.Getenv("EVM_PRIVATE_KEY"))
	if priv == "" {
		log.Println("DEBUG: EVM_PRIVATE_KEY is empty")
		return nil, fmt.Errorf("EVM_PRIVATE_KEY not configured")
	}

	privateKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		return nil, fmt.Errorf("Invalid EVM_PRIVATE_KEY: %v", err)
	}

	// 2. Get Chain ID from client directly (Chain-Agnostic)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		chainID = big.NewInt(80002) // Fallback
		if cid, err := strconv.ParseInt(os.Getenv("L2_CHAIN_ID"), 10, 64); err == nil {
			chainID = big.NewInt(cid)
		}
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// DYNAMIC GAS STRATEGY (EIP-1559)
	// We use aggressive GasTipCap (priority) and GasFeeCap (max) to ensure mining.
	if client != nil {
		head, err := client.HeaderByNumber(context.Background(), nil)
		if err == nil && head.BaseFee != nil {
			// 1. Suggest Priority Fee (Tip)
			tip, errTip := client.SuggestGasTipCap(context.Background())
			if errTip != nil {
				tip = big.NewInt(2500000000) // Fallback to 2.5 Gwei tip
			} else {
				// Bump tip by 100% to ensure priority
				tip = new(big.Int).Mul(tip, big.NewInt(2))
			}

			// 2. Calculate Max Fee
			// maxFee = (baseFee * 2) + tip
			maxFee := new(big.Int).Mul(head.BaseFee, big.NewInt(2))
			maxFee.Add(maxFee, tip)

			auth.GasTipCap = tip
			auth.GasFeeCap = maxFee
			auth.GasPrice = nil // Ensure we use EIP-1559
		} else {
			// Fallback to bumped Legacy GasPrice if EIP-1559 is unavailable
			gasPrice, errGas := client.SuggestGasPrice(context.Background())
			if errGas == nil {
				auth.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2)) // 2x legacy price
			}
		}
	}

	// GAS_LIMIT: Safe default for Sepolia (500k)
	auth.GasLimit = 500000 
	if gl := strings.TrimSpace(os.Getenv("GAS_LIMIT")); gl != "" {
		gl = strings.Trim(gl, `"'`)
		if glBig, ok := new(big.Int).SetString(gl, 10); ok && glBig.Sign() > 0 {
			auth.GasLimit = glBig.Uint64()
		}
	}

	return auth, nil
}

func ensureContractVerified(client *ethclient.Client, addr common.Address, label string) error {
	cacheKey := strings.ToLower(addr.Hex())
	if _, ok := verifiedContractAddrs.Load(cacheKey); ok {
		return nil
	}

	code, err := client.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return fmt.Errorf("failed to inspect %s contract: %w", label, err)
	}
	if len(code) == 0 {
		return fmt.Errorf("no contract code at %s address %s", label, addr.Hex())
	}

	verifiedContractAddrs.Store(cacheKey, true)
	return nil
}

func loadManifestoMap(ctx context.Context, electionAddress string) map[string]string {
	if candidateCollection == nil {
		return nil
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := candidateCollection.Find(queryCtx, bson.M{"electionAddress": electionAddress})
	if err != nil {
		log.Printf("loadManifestoMap: Find error for %s: %v", electionAddress, err)
		return nil
	}
	defer cursor.Close(queryCtx)

	var docs []CandidateDocument
	if err := cursor.All(queryCtx, &docs); err != nil {
		log.Printf("loadManifestoMap: decode error for %s: %v", electionAddress, err)
		return nil
	}

	manifestoMap := make(map[string]string, len(docs))
	for _, doc := range docs {
		if doc.Email != "" && doc.ManifestoUrl != "" {
			manifestoMap[doc.Email] = doc.ManifestoUrl
		}
	}
	return manifestoMap
}

// Global Nonce Manager for High Concurrency
var (
	nonceMutex sync.Mutex
	nonces     = make(map[uint64]uint64)
	nonceInits = make(map[uint64]bool)
)

func resyncNonce(client *ethclient.Client, address common.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID for resync: %v", err)
	}
	cid := chainID.Uint64()

	pendingNonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		return err
	}

	nonceMutex.Lock()
	defer nonceMutex.Unlock()
	nonces[cid] = pendingNonce
	nonceInits[cid] = true
	log.Printf("[NONCE] Resynced nonce for chain %d: %d", cid, pendingNonce)
	return nil
}

// getNextNonce guarantees a strictly increasing nonce per chain for the admin wallet.
func getNextNonce(client *ethclient.Client, address common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}
	cid := chainID.Uint64()

	nonceMutex.Lock()
	if !nonceInits[cid] {
		nonceMutex.Unlock()
		if err := resyncNonce(client, address); err != nil {
			return nil, fmt.Errorf("initialize nonce for %s on chain %d: %w", address.Hex(), cid, err)
		}
		nonceMutex.Lock()
	}

	next := nonces[cid]
	nonces[cid]++
	nonceMutex.Unlock()

	return new(big.Int).SetUint64(next), nil
}

func isNonceError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nonce too low") ||
		strings.Contains(msg, "nonce too high") ||
		strings.Contains(msg, "replacement transaction underpriced") ||
		strings.Contains(msg, "already known")
}

func submitChainTx(
	client *ethclient.Client,
	makeAuth func() (*bind.TransactOpts, error),
	submit func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {
	// 1. Get chain ID for logging
	chainID, _ := client.ChainID(context.Background())
	cid := uint64(0)
	if chainID != nil {
		cid = chainID.Uint64()
	}

	auth, err := makeAuth()
	if err != nil {
		return nil, err
	}

	// 2. Fetch Nonce
	nonce, err := getNextNonce(client, auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce = nonce

	log.Printf("[BC-TX] Submitting tx to chain %d with nonce %d for %s", cid, nonce.Uint64(), auth.From.Hex())

	// 3. Submit
	tx, err := submit(auth)
	if err == nil {
		log.Printf("[BC-TX] Success! Tx sent: %s (Chain %d, Nonce %d)", tx.Hash().Hex(), cid, nonce.Uint64())
		return tx, nil
	}

	// 4. RECOVERY: If ANY error occurred (Nonce or Revert), we MUST resync the nonce.
	log.Printf("[RECOVERY] Tx failed on chain %d (Nonce %d): %v. Resyncing nonce to prevent deadlock.", cid, nonce.Uint64(), err)
	_ = resyncNonce(client, auth.From)

	if !isNonceError(err) {
		return nil, err
	}

	// 5. NONCE SPECIFIC RETRY
	log.Printf("[WARN] Nonce error on chain %d: %v. Retrying with fresh resynced nonce...", cid, err)
	if resyncErr := resyncNonce(client, auth.From); resyncErr != nil {
		return nil, fmt.Errorf("retry aborted: nonce resync failed: %w", resyncErr)
	}

	authRetry, authErr := makeAuth()
	if authErr != nil {
		return nil, authErr
	}
	nonceRetry, nonceErr := getNextNonce(client, authRetry.From)
	if nonceErr != nil {
		return nil, nonceErr
	}
	authRetry.Nonce = nonceRetry

	log.Printf("[BC-TX-RETRY] Retrying on chain %d with nonce %d", cid, nonceRetry.Uint64())
	txRetry, retryErr := submit(authRetry)
	if retryErr != nil {
		log.Printf("[BC-TX-RETRY] Failed again: %v", retryErr)
		_ = resyncNonce(client, authRetry.From) // resync again on final failure
		return nil, retryErr
	}

	log.Printf("[BC-TX-RETRY] Success on retry! Tx: %s", txRetry.Hash().Hex())
	return txRetry, nil
}

// normalizeFactoryAddr returns a validated, 0x-prefixed factory address string and the parsed common.Address.
func normalizeFactoryAddr() (string, common.Address, error) {
	raw := strings.TrimSpace(os.Getenv("L2_FACTORY_CONTRACT_ADDRESS"))
	raw = strings.Trim(raw, `"'`)
	if raw == "" {
		return "", common.Address{}, fmt.Errorf("L2_FACTORY_CONTRACT_ADDRESS not set")
	}
	// if address is 40 hex chars without 0x, add prefix
	if len(raw) == 40 && !strings.HasPrefix(raw, "0x") {
		raw = "0x" + raw
	}
	if !common.IsHexAddress(raw) {
		return raw, common.Address{}, fmt.Errorf("L2_FACTORY_CONTRACT_ADDRESS is not a valid hex address: %q", raw)
	}
	return raw, common.HexToAddress(raw), nil
}

// normalizeAddrParam ensures we get a full 0x-prefixed address (accepts 40-char without 0x).
// Returns normalized string or error if clearly invalid.
func normalizeAddrParam(param string) (string, error) {
	s := strings.TrimSpace(param)
	s = strings.Trim(s, `"'`)
	if s == "" {
		return "", fmt.Errorf("empty address")
	}
	// If looks like 40 hex chars without 0x
	if len(s) == 40 && !strings.HasPrefix(s, "0x") {
		s = "0x" + s
	}
	// If it's shorter than minimal plausible address (0x + 6 hex etc) treat as truncated
	if strings.HasPrefix(s, "0x") && len(s) < 10 { // arbitrary small threshold
		return "", fmt.Errorf("address too short / truncated")
	}
	if !common.IsHexAddress(s) {
		return "", fmt.Errorf("invalid hex address")
	}
	return s, nil
}

// -- CREATE ELECTION --
func CreateElection(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)

	var req struct {
		CompanyEmail        string `json:"company_email"`
		ElectionName        string `json:"election_name"`
		ElectionDescription string `json:"election_description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CompanyEmail == "" || req.ElectionName == "" || req.ElectionDescription == "" {
		respondError(w, http.StatusBadRequest, "company_email, election_name and election_description are required")
		return
	}

	client, err := getClient()
	if err != nil {
		log.Printf("CreateElection: getClient error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to connect to ethereum node")
		return
	}

	// Validate factory address early
	factoryRaw, factoryAddr, err := normalizeFactoryAddr()
	if err != nil {
		log.Printf("CreateElection: invalid L2_FACTORY_CONTRACT_ADDRESS: %v", err)
		respondError(w, http.StatusInternalServerError, "L2_FACTORY_CONTRACT_ADDRESS not set or invalid")
		return
	}

	if err := ensureContractVerified(client, factoryAddr, "factory"); err != nil {
		log.Printf("CreateElection: factory verification error for %s: %v", factoryRaw, err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Setup binding
	factory, err := bindings.NewElectionFact(factoryAddr, client)
	if err != nil {
		log.Printf("CreateElection: factory binding error for %s: %v", factoryRaw, err)
		respondError(w, http.StatusInternalServerError, "failed to bind to factory contract")
		return
	}

	// PROACTIVE NONCE RESYNC: Ensure we have the latest state before doing a critical admin task
	adminAuth, adminErr := getAuth(client)
	if adminErr == nil {
		_ = resyncNonce(client, adminAuth.From)
	}

	// Submit CreateElection tx
	tx, err := submitChainTx(
		client,
		func() (*bind.TransactOpts, error) {
			return getAuth(client)
		},
		func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return factory.CreateElection(auth, req.CompanyEmail, req.ElectionName, req.ElectionDescription)
		},
	)
	if err != nil {
		log.Printf("CreateElection: transact failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create election on blockchain: "+err.Error())
		return
	}

	// Respond immediately that transaction was submitted
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":           "success",
		"message":          "createElection transaction submitted",
		"election_address": "",
		"confirmed":        false,
		"data": map[string]interface{}{
			"txHash": tx.Hash().Hex(),
		},
	})

	// Wait for mining and fetch deployed address asynchronously
	go func() {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 600*time.Second) // generous 10 Min timeout
		defer cancel2()
		receipt, werr := bind.WaitMined(ctx2, client, tx)
		if werr != nil {
			log.Printf("[ALCHEMY] CreateElection: WaitMined error: %v", werr)
			return
		}
		if receipt.Status != 1 {
			log.Printf("[ALCHEMY] CreateElection: receipt indicates revert tx %s", tx.Hash().Hex())
			return
		}

		// After tx mined (or if not waiting), try to read deployed address via factory's GetDeployedElection
		// Note: We need a fresh context here since r.Context() might be cancelled when the HTTP request ends
		callCtx, callCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer callCancel()

		factoryCaller, err := bindings.NewElectionFactCaller(factoryAddr, client)
		if err != nil {
			// not fatal; we'll still return tx hash
			log.Printf("CreateElection: NewElectionFactoryCaller binding error for %s: %v", factoryRaw, err)
		} else {
			callOpts := &bind.CallOpts{Pending: false, Context: callCtx}

			// Declare vars to be used after the block
			var deployedAddr common.Address
			var name, desc string

			elections, derr := factoryCaller.GetDeployedElections(callOpts, req.CompanyEmail)
			if derr != nil {
				log.Printf("CreateElection: GetDeployedElections read error for email %s: %v", req.CompanyEmail, derr)
			} else if len(elections) > 0 {
				latest := elections[len(elections)-1]
				deployedAddr = latest.DeployedAddress
				name = latest.ElN
				desc = latest.ElD
			} else {
				// No elections found
				log.Printf("CreateElection: factory returned empty election list for email %s after create tx", req.CompanyEmail)
			}

			if deployedAddr != (common.Address{}) {
				addrHex := deployedAddr.Hex()
				log.Printf("[ALCHEMY] CreateElection async success. Deployed at: %s", addrHex)
				// AUDIT LOG
				go LogAction(addrHex, "ELECTION_CREATED", req.CompanyEmail, fmt.Sprintf("Created election '%s'", req.ElectionName))
				// METADATA INIT
				go EnsureMetadata(addrHex, name, desc)
			} else {
				log.Printf("CreateElection async: factory returned zero address for email %s after create tx", req.CompanyEmail)
			}
		}
	}()
}

// VoteCandidate uses the election binding to cast a vote.
func VoteCandidate(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)

	var req struct {
		ElectionAddress string `json:"election_address"`
		CandidateID     int64  `json:"candidate_id"`
		VoterEmail      string `json:"voter_email"`
		OTP             string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.ElectionAddress) == "" || req.VoterEmail == "" {
		respondError(w, http.StatusBadRequest, "election_address and voter_email are required")
		return
	}

	// normalize election address
	addrNorm, err := normalizeAddrParam(req.ElectionAddress)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid election_address: "+err.Error())
		return
	}

	// CHECK PHASES
	active, reason := IsElectionActive(addrNorm)
	if !active {
		respondError(w, http.StatusBadRequest, "Voting not allowed: "+reason)
		return
	}

	// CHECK VERIFICATION
	if verified := IsVoterVerified(req.VoterEmail, addrNorm); !verified {
		respondError(w, http.StatusForbidden, "Voter not verified. Please contact election admin.")
		return
	}

	// MFA CHECK (Commented out)
	/*
		if ok := VerifyAndDeleteOTP(req.VoterEmail, req.OTP); !ok {
			respondError(w, http.StatusUnauthorized, "Invalid or expired OTP")
			return
		}
	*/

	// 2. INSTANT DB CHECK: Check if voter already has "Voted" status in registrations
	if voterCollection != nil {
		var v Voter
		addrRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(addrNorm) + "$", "$options": "i"}
		emailRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(req.VoterEmail) + "$", "$options": "i"}
		
		err := voterCollection.FindOne(context.Background(), bson.M{
			"email": emailRegex,
			"registrations": bson.M{"$elemMatch": bson.M{
				"election_address": addrRegex,
				"status": bson.M{"$regex": "^voted$", "$options": "i"},
			}},
		}).Decode(&v)
		if err == nil {
			respondError(w, http.StatusBadRequest, "Double voting detected: You have already cast your vote for this election.")
			return
		}
	}

	// 3. JOB CHECK: Check if a vote is currently in our processing queue
	existingJob, _ := getVoteJobByVoter(addrNorm, req.VoterEmail)
	if existingJob != nil && existingJob.Status != "failed" {
		msg := "You have already voted in this election."
		if existingJob.Status == "queued" || existingJob.Status == "submitted" {
			msg = "Your vote is currently being processed on the blockchain. Please wait."
		}
		respondError(w, http.StatusBadRequest, msg)
		return
	}

	// IMMEDIATE LOCK: Mark as "Voted" in DB before queuing
	if voterCollection != nil {
		vCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		addrRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(addrNorm) + "$", "$options": "i"}
		emailRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(req.VoterEmail) + "$", "$options": "i"}

		_, _ = voterCollection.UpdateOne(vCtx,
			bson.M{"email": emailRegex, "registrations.election_address": addrRegex},
			bson.M{"$set": bson.M{"registrations.$.status": "Voted"}},
		)
	}

	job, err := enqueueVoteJob(addrNorm, req.CandidateID, req.VoterEmail)
	if err != nil {
		log.Printf("VoteCandidate: enqueue error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to queue vote transaction")
		return
	}

	message := "vote transaction queued for blockchain submission"
	if job.Status != "queued" {
		message = "vote transaction already queued"
	}

	respondJSON(w, http.StatusAccepted, BlockchainResponse{
		Status:  "success",
		Message: message,
		Data: map[string]interface{}{
			"jobId":       job.ID.Hex(),
			"txHash":      job.TxHash,
			"queueStatus": job.Status,
		},
	})
}

// GetElectionCandidates - improved and robust
func GetElectionCandidates(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)

	vars := mux.Vars(r)
	rawAddr := vars["address"]
	log.Printf("GetElectionCandidates: raw address param: %q\n", rawAddr)

	if rawAddr == "" {
		respondError(w, http.StatusBadRequest, "missing election address")
		return
	}

	// Sanitize: trim spaces and possible surrounding quotes
	addrStr := strings.TrimSpace(rawAddr)
	addrStr = strings.Trim(addrStr, `"'`)

	// Treat common "no value" strings as empty -> db fallback
	if addrStr == "" || strings.EqualFold(addrStr, "null") || strings.EqualFold(addrStr, "undefined") {
		log.Printf("GetElectionCandidates: address param empty or null-like (%q) - using DB fallback\n", rawAddr)
		tryDBFallbackWithMessage(w, addrStr, "invalid or truncated election address")
		return
	}

	// If the input looks like a truncated hex (starts with 0x but length < 42),
	// attempt to resolve it by searching the DB for an electionAddress that starts with this prefix.
	if strings.HasPrefix(addrStr, "0x") && len(addrStr) < 42 {
		prefix := addrStr
		log.Printf("GetElectionCandidates: received truncated address prefix: %q - attempting DB prefix lookup\n", prefix)

		// If candidateCollection isn't set, we can't search DB - just fallback.
		if candidateCollection == nil {
			log.Printf("GetElectionCandidates: no candidateCollection available for prefix lookup; using DB fallback\n")
			tryDBFallbackWithMessage(w, addrStr, "invalid or truncated election address")
			return
		}

		regexPattern := "^" + regexp.QuoteMeta(prefix)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{"electionAddress": bson.M{"$regex": regexPattern, "$options": "i"}}
		findOpts := options.Find()
		findOpts.SetLimit(5)
		cursor, err := candidateCollection.Find(ctx, filter, findOpts)
		if err != nil {
			log.Printf("GetElectionCandidates: DB prefix search error: %v\n", err)
			tryDBFallbackWithMessage(w, addrStr, "db lookup failed while resolving truncated address")
			return
		}
		defer cursor.Close(ctx)

		found := map[string]struct{}{}
		for cursor.Next(ctx) {
			var doc CandidateDocument
			if err := cursor.Decode(&doc); err != nil {
				log.Printf("GetElectionCandidates: cursor decode error during prefix lookup: %v\n", err)
				continue
			}
			if doc.ElectionAddress != "" {
				found[doc.ElectionAddress] = struct{}{}
			}
		}
		if err := cursor.Err(); err != nil {
			log.Printf("GetElectionCandidates: cursor iteration error: %v\n", err)
		}

		if len(found) == 1 {
			var resolved string
			for k := range found {
				resolved = k
				break
			}
			log.Printf("GetElectionCandidates: resolved truncated prefix %q -> full address %s via DB\n", prefix, resolved)
			addrStr = resolved
		} else if len(found) > 1 {
			log.Printf("GetElectionCandidates: truncated prefix %q matched multiple addresses (%d) - returning ambiguous error\n", prefix, len(found))
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": "ambiguous truncated election identifier; multiple elections match this prefix - please provide the full address",
				"matches": len(found),
			})
			return
		} else {
			log.Printf("GetElectionCandidates: truncated prefix %q did not match any electionAddress in DB - using DB fallback\n", prefix)
			tryDBFallbackWithMessage(w, addrStr, "invalid or truncated election address")
			return
		}
	}

	// Accept 40-hex without 0x
	if len(addrStr) == 40 && !strings.HasPrefix(addrStr, "0x") {
		if common.IsHexAddress("0x" + addrStr) {
			addrStr = "0x" + addrStr
		}
	}

	// If it's still not a hex address, attempt to resolve as a company email via factory lookup.
	if !common.IsHexAddress(addrStr) {
		log.Printf("GetElectionCandidates: address param %q is not hex - trying factory lookup as email\n", addrStr)

		// Try to read factory address and call GetDeployedElection(email)
		factoryRaw, factoryAddr, ferr := normalizeFactoryAddr()
		if ferr != nil {
			log.Printf("GetElectionCandidates: cannot resolve non-hex param via factory: %v\n", ferr)
			tryDBFallbackWithMessage(w, addrStr, "invalid or truncated election address")
			return
		}

		client, cerr := getClient()
		if cerr != nil {
			log.Printf("GetElectionCandidates: getClient error while resolving email: %v\n", cerr)
			tryDBFallbackWithMessage(w, addrStr, "failed to connect to ethereum node while resolving email")
			return
		}

		if err := ensureContractVerified(client, factoryAddr, "factory"); err != nil {
			log.Printf("GetElectionCandidates: factory verification error for %s: %v\n", factoryRaw, err)
			tryDBFallbackWithMessage(w, addrStr, err.Error())
			return
		}

		factoryCaller, ferr2 := bindings.NewElectionFactCaller(factoryAddr, client)
		if ferr2 != nil {
			log.Printf("GetElectionCandidates: factory caller binding error for %s: %v\n", factoryRaw, ferr2)
			tryDBFallbackWithMessage(w, addrStr, "factory binding error while resolving email")
			return
		}

		callOpts := &bind.CallOpts{Context: r.Context(), Pending: false}
		elections, gerr := factoryCaller.GetDeployedElections(callOpts, addrStr)
		if gerr != nil {
			log.Printf("GetElectionCandidates: GetDeployedElections error for %q: %v\n", addrStr, gerr)
			tryDBFallbackWithMessage(w, addrStr, "factory lookup failed for provided identifier")
			return
		}
		if len(elections) == 0 {
			log.Printf("GetElectionCandidates: factory returned empty list for %q - falling back to DB\n", addrStr)
			tryDBFallbackWithMessage(w, addrStr, "no deployed election found for provided identifier")
			return
		}
		// Use latest
		deployedAddr := elections[len(elections)-1].DeployedAddress

		if deployedAddr == (common.Address{}) {
			log.Printf("GetElectionCandidates: factory returned zero address for %q - falling back to DB\n", addrStr)
			tryDBFallbackWithMessage(w, addrStr, "no deployed election found for provided identifier")
			return
		}
		// resolved - normalize to hex address and continue onchain flow
		addrStr = deployedAddr.Hex()
		log.Printf("GetElectionCandidates: resolved %q -> onchain address %s via factory %s\n", rawAddr, addrStr, factoryRaw)
	}

	// At this point addrStr should be a valid hex address (0x...)
	if !common.IsHexAddress(addrStr) {
		log.Printf("GetElectionCandidates: invalid election address after normalization: %q\n", addrStr)
		tryDBFallbackWithMessage(w, addrStr, "invalid or truncated election address")
		return
	}

	cacheID := strings.ToLower(addrStr)
	if cached, ok := getCachedValue(cacheKey("election", cacheID, "candidates")); ok {
		respondJSON(w, http.StatusOK, cached)
		return
	}

	// Connect to node
	client, err := getClient()
	if err != nil {
		log.Printf("GetElectionCandidates: getClient error: %v\n", err)
		tryDBFallbackWithMessage(w, addrStr, "failed to connect to ethereum node: "+err.Error())
		return
	}

	addr := common.HexToAddress(addrStr)
	if err := ensureContractVerified(client, addr, "election"); err != nil {
		log.Printf("GetElectionCandidates: contract verification error for %s: %v\n", addrStr, err)
		tryDBFallbackWithMessage(w, addrStr, err.Error())
		return
	}

	// Bind contract and read
	contract, err := bindings.NewElection(addr, client)
	if err != nil {
		log.Printf("GetElectionCandidates: bindings.NewElection error: %v\n", err)
		tryDBFallbackWithMessage(w, addrStr, "failed to bind contract: "+err.Error())
		return
	}

	callOpts := &bind.CallOpts{Context: r.Context(), Pending: false}
	numCandidates, err := contract.GetNumOfCandidates(callOpts)
	if err != nil {
		log.Printf("GetElectionCandidates: GetNumOfCandidates error for %s: %v\n", addrStr, err)
		tryDBFallbackWithMessage(w, addrStr, "failed to fetch candidate count from chain: "+err.Error())
		return
	}

	n := numCandidates.Int64()
	manifestoMap := loadManifestoMap(r.Context(), addrStr)
	candidates := make([]Candidate, 0, n)
	for i := int64(0); i < n; i++ {
		name, desc, imgHash, voteCount, email, err := contract.GetCandidate(callOpts, big.NewInt(i))
		if err != nil {
			log.Printf("GetElectionCandidates: GetCandidate(%d) error for %s: %v\n", i, addrStr, err)
			tryDBFallbackWithMessage(w, addrStr, fmt.Sprintf("failed to fetch candidate %d from chain: %v", i, err))
			return
		}
		manifesto := ""
		if manifestoMap != nil {
			manifesto = manifestoMap[email]
		}

		candidates = append(candidates, Candidate{
			Name:         name,
			Description:  desc,
			ImageHash:    imgHash,
			VoteCount:    voteCount,
			Email:        email,
			ManifestoUrl: manifesto,
		})
	}

	// Return success with source = "onchain"
	log.Printf("[SUCCESS] Successfully fetched %d candidates from blockchain for %s\n", len(candidates), addrStr)
	payload := map[string]interface{}{
		"status":     "success",
		"source":     "onchain",
		"candidates": candidates,
	}
	setCachedValue(cacheKey("election", cacheID, "candidates"), payload, 15*time.Second)
	respondJSON(w, http.StatusOK, payload)
}

// tryDBFallbackWithMessage returns DB candidates and includes the provided message in result.detail
func tryDBFallbackWithMessage(w http.ResponseWriter, electionAddress, detail string) {
	// Attempt to return DB candidates to keep UI usable
	// NOTE: candidateCollection should be initialized elsewhere in your app
	if candidateCollection == nil {
		// no DB available - return error JSON
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "no db fallback available",
			"detail":  detail,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	filter := bson.M{"electionAddress": electionAddress}
	cursor, err := candidateCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("tryDBFallbackWithMessage: Find error: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "db lookup failed",
			"detail":  detail + " | db find error: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var docs []CandidateDocument
	if err := cursor.All(ctx, &docs); err != nil {
		log.Printf("tryDBFallbackWithMessage: cursor.All error: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "db decode failed",
			"detail":  detail + " | db decode error: " + err.Error(),
		})
		return
	}

	// Map docs to lightweight candidate shape
	candidates := make([]map[string]interface{}, 0, len(docs))
	for _, d := range docs {
		candidates = append(candidates, map[string]interface{}{
			"name":         d.Name,
			"description":  d.Description,
			"imageHash":    d.ImageHash,
			"email":        d.Email,
			"manifestoUrl": d.ManifestoUrl,
			"txHash":       d.TxHash,
			"createdAt":    d.CreatedAt,
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":     "success",
		"source":     "db_fallback",
		"detail":     detail,
		"candidates": candidates,
	})
}

func GetElectionInfo(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)

	vars := mux.Vars(r)
	rawAddr := strings.TrimSpace(vars["address"])
	rawAddr = strings.Trim(rawAddr, `"'`)
	if rawAddr == "" {
		respondError(w, http.StatusBadRequest, "election address parameter is required")
		return
	}

	// Accept 40-char hex without 0x for convenience
	if len(rawAddr) == 40 && !strings.HasPrefix(rawAddr, "0x") {
		rawAddr = "0x" + rawAddr
	}

	// If not hex, try to resolve as email via factory
	if !common.IsHexAddress(rawAddr) {
		// normalize factory
		_, factoryAddr, ferr := normalizeFactoryAddr()
		if ferr == nil {
			client, cerr := getClient()
			if cerr == nil {
				factoryCaller, ferr2 := bindings.NewElectionFactCaller(factoryAddr, client)
				if ferr2 == nil {
					callOpts := &bind.CallOpts{Context: r.Context(), Pending: false}
					elections, gerr := factoryCaller.GetDeployedElections(callOpts, rawAddr)
					if gerr == nil && len(elections) > 0 {
						// Use latest
						latest := elections[len(elections)-1]
						if latest.DeployedAddress != (common.Address{}) {
							rawAddr = latest.DeployedAddress.Hex()
						}
					}
				}
			}
		}
	}

	if !common.IsHexAddress(rawAddr) {
		respondError(w, http.StatusBadRequest, "invalid election address or unresolved email")
		return
	}
	_ = common.HexToAddress(rawAddr)

	cacheID := strings.ToLower(rawAddr)
	if cached, ok := getCachedValue(cacheKey("election", cacheID, "info")); ok {
		respondJSON(w, http.StatusOK, cached)
		return
	}

	// Use MongoDB for all dashboard stats with case-insensitive address matching
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Robust case-insensitive lookup to handle checksummed vs lowercase address mismatches
	addrRegex := bson.M{"$regex": "^" + regexp.QuoteMeta(rawAddr) + "$", "$options": "i"}

	// Use an extremely robust $or filter to handle potential schema inconsistencies
	// (e.g. array of objects vs array of strings, snake_case vs camelCase)
	votersFilter := bson.M{"$or": []bson.M{
		{"registrations.election_address": addrRegex},
		{"registrations.electionAddress": addrRegex},
		{"registrations": addrRegex},
		{"election_address": addrRegex},
		{"electionAddress": addrRegex},
	}}

	votersCount, _ := voterCollection.CountDocuments(ctx, votersFilter)
	candidatesCount, _ := candidateCollection.CountDocuments(ctx, bson.M{"electionAddress": addrRegex})

	var meta ElectionMetadata
	_ = metadataCollection.FindOne(ctx, bson.M{"election_address": addrRegex}).Decode(&meta)

	payload := BlockchainResponse{
		Status:  "success",
		Message: "election info retrieved",
		Data: map[string]interface{}{
			"voters_count":     fmt.Sprintf("%d", votersCount),
			"candidates_count": fmt.Sprintf("%d", candidatesCount),
			"election_name":    meta.ElectionName,
			"election_desc":    meta.ElectionDesc,
			"election_addr":    rawAddr,
		},
	}
	setCachedValue(cacheKey("election", cacheID, "info"), payload, 15*time.Second)
	respondJSON(w, http.StatusOK, payload)
}

// UploadImage (disabled)
func UploadImage(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	respondJSON(w, http.StatusNotFound, BlockchainResponse{
		Status:  "error",
		Message: "image upload endpoint has been removed; image uploads are disabled",
	})
}
