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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
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

// getClient connects to L2_NODE_URL (with timeout)
func getClient() (*ethclient.Client, error) {
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
	return client, nil
}

// getAuth creates a transact opts using EVM_PRIVATE_KEY and L2_CHAIN_ID.
// It requires a valid ethclient to fetch the current block nonce safely.
func getAuth() (*bind.TransactOpts, error) {
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

	// 2. Get Chain ID
	chainIDStr := strings.TrimSpace(os.Getenv("L2_CHAIN_ID"))
	chainIDInt, err := strconv.ParseInt(chainIDStr, 10, 64)
	if err != nil || chainIDInt == 0 {
		chainIDInt = 80002 // Default to Amoy
	}
	chainID := big.NewInt(chainIDInt)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// optional GAS_LIMIT override (env expects decimal integer)
	if gl := strings.TrimSpace(os.Getenv("GAS_LIMIT")); gl != "" {
		gl = strings.Trim(gl, `"'`)
		if glBig, ok := new(big.Int).SetString(gl, 10); ok && glBig.Sign() > 0 {
			auth.GasLimit = glBig.Uint64()
		}
	}

	// optional GAS_PRICE override (env expects decimal wei)
	if gp := strings.TrimSpace(os.Getenv("GAS_PRICE")); gp != "" {
		gp = strings.Trim(gp, `"'`)
		if gpBig, ok := new(big.Int).SetString(gp, 10); ok && gpBig.Sign() > 0 {
			auth.GasPrice = gpBig
		}
	}

	return auth, nil
}

// Global Nonce Manager for High Concurrency
var (
	nonceMutex sync.Mutex
	lastNonce  uint64
)

// getNextNonce guarantees a strictly increasing nonce for the admin wallet, even during extreme concurrency.
func getNextNonce(client *ethclient.Client, address common.Address) *big.Int {
	nonceMutex.Lock()
	defer nonceMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch pending nonce from the network
	pendingNonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		log.Printf("[WARN] Failed to fetch pending nonce for %s: %v. Using last known nonce.", address.Hex(), err)
		pendingNonce = lastNonce
	}

	// If our internal counter is higher (meaning we have unbroadcasted local txs), use it.
	if lastNonce >= pendingNonce {
		pendingNonce = lastNonce + 1
	}

	lastNonce = pendingNonce
	return new(big.Int).SetUint64(pendingNonce)
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
	defer client.Close()

	auth, err := getAuth()
	if err != nil {
		log.Printf("CreateElection: getAuth error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create transaction signer")
		return
	}
	auth.Nonce = getNextNonce(client, auth.From) // Set nonce here after client is available

	// Validate factory address early
	factoryRaw, factoryAddr, err := normalizeFactoryAddr()
	if err != nil {
		log.Printf("CreateElection: invalid L2_FACTORY_CONTRACT_ADDRESS: %v", err)
		respondError(w, http.StatusInternalServerError, "L2_FACTORY_CONTRACT_ADDRESS not set or invalid")
		return
	}

	// quick: check code at factory address
	code, cerr := client.CodeAt(context.Background(), factoryAddr, nil)
	if cerr != nil {
		log.Printf("CreateElection: CodeAt error for %s: %v", factoryRaw, cerr)
		respondError(w, http.StatusInternalServerError, "failed to inspect factory contract")
		return
	}
	if len(code) == 0 {
		log.Printf("CreateElection: no contract code at factory address %s", factoryRaw)
		respondError(w, http.StatusInternalServerError, "no contract code at configured L2_FACTORY_CONTRACT_ADDRESS")
		return
	}

	// Setup binding
	factory, err := bindings.NewElectionFact(factoryAddr, client)
	if err != nil {
		log.Printf("CreateElection: factory binding error for %s: %v", factoryRaw, err)
		respondError(w, http.StatusInternalServerError, "failed to bind to factory contract")
		return
	}

	// Submit CreateElection tx
	tx, err := factory.CreateElection(auth, req.CompanyEmail, req.ElectionName, req.ElectionDescription)
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
		ctx2, cancel2 := context.WithTimeout(context.Background(), 120*time.Second) // generous timeout
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

	// MFA CHECK
	if ok := VerifyAndDeleteOTP(req.VoterEmail, req.OTP); !ok {
		respondError(w, http.StatusUnauthorized, "Invalid or expired OTP")
		return
	}

	client, err := getClient()
	if err != nil {
		log.Printf("VoteCandidate: getClient error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to connect to ethereum node")
		return
	}
	defer client.Close()

	auth, err := getAuth()
	if err != nil {
		log.Printf("VoteCandidate: getAuth error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create transaction signer")
		return
	}
	auth.Nonce = getNextNonce(client, auth.From)

	contractAddr := common.HexToAddress(addrNorm)
	// check contract code present
	code, cerr := client.CodeAt(context.Background(), contractAddr, nil)
	if cerr != nil {
		log.Printf("VoteCandidate: CodeAt error for %s: %v", addrNorm, cerr)
		respondError(w, http.StatusInternalServerError, "failed to inspect contract code")
		return
	}
	if len(code) == 0 {
		log.Printf("VoteCandidate: no contract code at address %s", addrNorm)
		respondError(w, http.StatusBadRequest, "no contract code at given election address")
		return
	}

	contract, err := bindings.NewElection(contractAddr, client)
	if err != nil {
		log.Printf("VoteCandidate: binding error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to bind to election contract")
		return
	}

	tx, err := contract.Vote(auth, big.NewInt(req.CandidateID), req.VoterEmail)
	if err != nil {
		log.Printf("VoteCandidate: vote transact error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to submit vote transaction: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, BlockchainResponse{
		Status:  "success",
		Message: "vote transaction submitted to the blockchain",
		Data:    map[string]interface{}{"txHash": tx.Hash().Hex()},
	})

	// Wait for mining asynchronously
	go func() {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 120*time.Second) // 2 Min timeout
		defer cancel2()
		receipt, werr := bind.WaitMined(ctx2, client, tx)
		if werr != nil {
			log.Printf("[ALCHEMY] Vote wait error: %v", werr)
		} else if receipt.Status != 1 {
			log.Printf("[ALCHEMY] Vote transaction reverted for tx %s", tx.Hash().Hex())
		} else {
			log.Printf("[ALCHEMY] Vote mined successfully in block %v", receipt.BlockNumber)
			// AUDIT LOG
			go LogAction(addrNorm, "VOTE_CAST", req.VoterEmail, "Voted successfully (mined)")
		}
	}()
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
		defer client.Close()

		// quick check factory code presence too
		fcode, ferrC := client.CodeAt(context.Background(), factoryAddr, nil)
		if ferrC != nil {
			log.Printf("GetElectionCandidates: CodeAt error for factory %s: %v\n", factoryRaw, ferrC)
			tryDBFallbackWithMessage(w, addrStr, "failed to inspect factory contract while resolving email")
			return
		}
		if len(fcode) == 0 {
			log.Printf("GetElectionCandidates: no contract code at factory address %s\n", factoryRaw)
			tryDBFallbackWithMessage(w, addrStr, "no factory contract code at configured address")
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

	// Connect to node
	client, err := getClient()
	if err != nil {
		log.Printf("GetElectionCandidates: getClient error: %v\n", err)
		tryDBFallbackWithMessage(w, addrStr, "failed to connect to ethereum node: "+err.Error())
		return
	}
	defer client.Close()

	// Check whether there is code at this address (if none -> no contract deployed)
	addr := common.HexToAddress(addrStr)
	code, err := client.CodeAt(r.Context(), addr, nil)
	if err != nil {
		log.Printf("GetElectionCandidates: CodeAt error for %s: %v\n", addrStr, err)
		tryDBFallbackWithMessage(w, addrStr, "failed to inspect contract code: "+err.Error())
		return
	}
	if len(code) == 0 {
		log.Printf("GetElectionCandidates: no contract code at address %s\n", addrStr)
		tryDBFallbackWithMessage(w, addrStr, "no contract code at given address")
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
	candidates := make([]Candidate, 0, n)
	for i := int64(0); i < n; i++ {
		name, desc, imgHash, voteCount, email, err := contract.GetCandidate(callOpts, big.NewInt(i))
		if err != nil {
			log.Printf("GetElectionCandidates: GetCandidate(%d) error for %s: %v\n", i, addrStr, err)
			tryDBFallbackWithMessage(w, addrStr, fmt.Sprintf("failed to fetch candidate %d from chain: %v", i, err))
			return
		}
		// merge manifesto from DB (inefficient loop but safe)
		manifesto := ""
		if candidateCollection != nil {
			var doc CandidateDocument
			if err := candidateCollection.FindOne(r.Context(), bson.M{"email": email, "electionAddress": addrStr}).Decode(&doc); err != nil {
				if err != mongo.ErrNoDocuments {
					log.Printf("GetElectionCandidates: manifesto lookup warning for %s: %v", email, err)
				}
			} else {
				manifesto = doc.ManifestoUrl
			}
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
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":     "success",
		"source":     "onchain",
		"candidates": candidates,
	})
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
				defer client.Close()
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

	respondJSON(w, http.StatusOK, BlockchainResponse{
		Status:  "success",
		Message: "election info retrieved",
		Data: map[string]interface{}{
			"voters_count":     fmt.Sprintf("%d", votersCount),
			"candidates_count": fmt.Sprintf("%d", candidatesCount),
			"election_name":    meta.ElectionName,
			"election_desc":    meta.ElectionDesc,
			"election_addr":    rawAddr,
		},
	})
}

// UploadImage (disabled)
func UploadImage(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	respondJSON(w, http.StatusNotFound, BlockchainResponse{
		Status:  "error",
		Message: "image upload endpoint has been removed; image uploads are disabled",
	})
}
