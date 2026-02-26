package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"math/big"
	"os"
	"strings"

	"MAJOR-PROJECT/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ElectionMetadata stores off-chain details about an election, primarily for scheduling/phases.
type ElectionMetadata struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ElectionAddress string             `bson:"election_address" json:"election_address"`

	// Dates for Phases
	StartDate time.Time `bson:"start_date" json:"start_date"`
	EndDate   time.Time `bson:"end_date" json:"end_date"`

	// Status derived from time or manual override
	Status string `bson:"status" json:"status"` // "UPCOMING", "ONGOING", "ENDED"

	// Display Details
	ElectionName string `bson:"election_name" json:"election_name"`
	ElectionDesc string `bson:"election_desc" json:"election_desc"`
}

var metadataCollection *mongo.Collection

// InitMetadataCollection initializes the collection
func InitMetadataCollection(client *mongo.Client, dbName string) {
	metadataCollection = client.Database(dbName).Collection("election_metadata")
	fmt.Println("[OK] Initialized election_metadata collection")
}

// EnsureMetadata creates a default metadata entry if one doesn't exist, and stores name/desc.
func EnsureMetadata(electionAddr, name, desc string) {
	if metadataCollection == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var meta ElectionMetadata
	err := metadataCollection.FindOne(ctx, bson.M{"election_address": electionAddr}).Decode(&meta)

	now := time.Now().UTC()
	defaultDuration := 7 * 24 * time.Hour // 7 days

	switch err {
	case mongo.ErrNoDocuments:
		// Create new
		newMeta := ElectionMetadata{
			ElectionAddress: electionAddr,
			StartDate:       now.Add(-1 * time.Hour), // Start slightly in past to avoid timezone races
			EndDate:         now.Add(defaultDuration),
			Status:          "ONGOING",
			ElectionName:    name,
			ElectionDesc:    desc,
		}
		metadataCollection.InsertOne(ctx, newMeta)
		fmt.Printf("[OK] Created metadata for %s (Expires: %s)\n", electionAddr, newMeta.EndDate)

	case nil:
		// Exists - Check if expired
		if now.After(meta.EndDate) && meta.Status != "ENDED" {
			fmt.Printf("[INFO] Election %s is past EndDate (%s)\n", electionAddr, meta.EndDate)
		}

		// Update Name/Desc if missing and provided
		if (meta.ElectionName == "" && name != "") || (meta.ElectionDesc == "" && desc != "") {
			update := bson.M{"$set": bson.M{}}
			if name != "" {
				update["$set"].(bson.M)["election_name"] = name
			}
			if desc != "" {
				update["$set"].(bson.M)["election_desc"] = desc
			}
			metadataCollection.UpdateOne(ctx, bson.M{"election_address": electionAddr}, update)
			fmt.Printf("[OK] Updated metadata details for %s\n", electionAddr)
		}
	}
}

// IsElectionActive checks if the current time is within the start/end window
func IsElectionActive(electionAddr string) (bool, string) {
	if metadataCollection == nil {
		return true, ""
	} // fallback if DB issue: allow voting

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var meta ElectionMetadata
	err := metadataCollection.FindOne(ctx, bson.M{"election_address": electionAddr}).Decode(&meta)
	if err != nil {
		// if no metadata found, assume it's legacy or open
		return true, ""
	}

	now := time.Now().UTC()
	if now.Before(meta.StartDate) {
		return false, fmt.Sprintf("Election has not started yet. Starts at %s UTC", meta.StartDate.Format("2006-01-02 15:04"))
	}
	if now.After(meta.EndDate) || meta.Status == "ENDED" {
		return false, "Election has ended."
	}
	return true, ""
}

// SetElectionDates Endpoint
func SetElectionDates(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ElectionAddress string `json:"election_address"`
		StartStr        string `json:"start_date"` // Expect RFC3339 or "2006-01-02T15:04"
		EndStr          string `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Parse times
	// Try a few layouts
	layout1 := "2006-01-02T15:04"

	start, err := time.Parse(time.RFC3339, req.StartStr)
	if err != nil {
		start, err = time.Parse(layout1, req.StartStr)
	}
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start_date format (use RFC3339 or YYYY-MM-DDTHH:MM)")
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndStr)
	if err != nil {
		end, err = time.Parse(layout1, req.EndStr)
	}
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid end_date format")
		return
	}

	if end.Before(start) {
		respondError(w, http.StatusBadRequest, "End date cannot be before start date")
		return
	}

	// Upsert to DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"election_address": req.ElectionAddress}
	update := bson.M{
		"$set": bson.M{
			"start_date": start,
			"end_date":   end,
			"status":     "SCHEDULED", // You might want logic to auto-calc status but this is fine
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err = metadataCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update dates")
		return
	}

	// Log it
	go LogAction(req.ElectionAddress, "SCHEDULE_UPDATE", "Admin", fmt.Sprintf("Dates updated: %s to %s", start, end))

	respondJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Election dates updated"})
}

// GetElectionMetadata Endpoint
func GetElectionMetadata(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	vars := mux.Vars(r)
	addr := vars["address"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var meta ElectionMetadata
	err := metadataCollection.FindOne(ctx, bson.M{"election_address": addr}).Decode(&meta)
	if err != nil {
		// Return 404 or just default
		respondError(w, http.StatusNotFound, "Metadata not found (election might use default open dates)")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   meta,
	})
}

// EndElection immediately stops an election
func EndElection(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	addr := vars["address"]
	if addr == "" {
		respondError(w, http.StatusBadRequest, "Address required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update Status=ENDED and EndDate=Now
	update := bson.M{
		"$set": bson.M{
			"status":   "ENDED",
			"end_date": time.Now().UTC(),
		},
	}
	_, err := metadataCollection.UpdateOne(ctx, bson.M{"election_address": addr}, update)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to end election")
		return
	}

	// --- L2 -> L1 ANCHORING LOGIC ---
	// Start anchoring asynchronously to not block the API response
	go func(electionAddress string) {
		log.Printf("[ANCHOR] Starting archiving process for election %s", electionAddress)

		// 1. Connect to L2 Client to read the results
		l2Url := strings.TrimSpace(os.Getenv("L2_NODE_URL"))
		if l2Url == "" {
			log.Println("[ANCHOR ERROR] L2_NODE_URL not set")
			return
		}

		l2Client, err := ethclient.Dial(l2Url)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to connect to L2: %v", err)
			return
		}
		defer l2Client.Close()

		l2Election, err := bindings.NewElection(common.HexToAddress(electionAddress), l2Client)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to bind L2 Election contract: %v", err)
			return
		}

		callOpts := &bind.CallOpts{Context: context.Background(), Pending: false}

		// 2. Fetch Winners and Metadata from L2
		title, _, err := l2Election.GetElectionDetails(callOpts)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to get title: %v", err)
			return
		}

		numVoters, err := l2Election.GetNumOfVoters(callOpts)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to get Total Voters: %v", err)
			return
		}

		winnerId, err := l2Election.WinnerCandidate(callOpts)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to get Winner ID: %v", err)
			return
		}

		winnerName, _, _, winningVotes, _, err := l2Election.GetCandidate(callOpts, winnerId)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to get Winner Details: %v", err)
			return
		}

		log.Printf("[ANCHOR] Results from L2: %s won '%s' with %s votes out of %s total voters", winnerName, title, winningVotes.String(), numVoters.String())

		// 3. Connect to L1 Client to archive the results
		l1Url := strings.TrimSpace(os.Getenv("L1_NODE_URL"))
		l1ArchiveAddr := strings.TrimSpace(os.Getenv("L1_ARCHIVE_CONTRACT_ADDRESS"))
		if l1Url == "" || l1ArchiveAddr == "" {
			log.Println("[ANCHOR ERROR] L1 config missing")
			return
		}

		l1Client, err := ethclient.Dial(l1Url)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to connect to L1: %v", err)
			return
		}
		defer l1Client.Close()

		l1Archive, err := bindings.NewBindings(common.HexToAddress(l1ArchiveAddr), l1Client)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to bind L1 Archive contract: %v", err)
			return
		}

		// 4. Create L1 Transactor
		privKeyStr := os.Getenv("EVM_PRIVATE_KEY")
		if privKeyStr == "" {
			log.Println("[ANCHOR ERROR] EVM_PRIVATE_KEY missing")
			return
		}

		privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyStr, "0x"))
		if err != nil {
			log.Printf("[ANCHOR ERROR] Invalid Private Key: %v", err)
			return
		}

		l1ChainIDVal, _ := new(big.Int).SetString(os.Getenv("L1_CHAIN_ID"), 10)
		if l1ChainIDVal == nil || l1ChainIDVal.Uint64() == 0 {
			l1ChainIDVal = big.NewInt(11155111)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(privKey, l1ChainIDVal)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to create L1 transactor: %v", err)
			return
		}

		// Fetch L1 Nonce specifically
		nonce, err := l1Client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			log.Printf("[ANCHOR ERROR] Failed to get L1 nonce: %v", err)
			return
		}
		auth.Nonce = big.NewInt(int64(nonce))

		// 5. Submit to L1
		tx, err := l1Archive.ArchiveResult(auth, common.HexToAddress(electionAddress), title, winnerName, winningVotes, numVoters)
		if err != nil {
			log.Printf("[ANCHOR ERROR] ArchiveResult tx failed: %v", err)
			return
		}

		log.Printf("[ANCHOR SUCCESS] Result for %s sent to L1 Sepolia at tx: %s", electionAddress, tx.Hash().Hex())

		// Log the anchoring completion in MongoDB audit
		go LogAction(electionAddress, "L1_ANCHOR_SUBMITTED", "System", fmt.Sprintf("Archived results to L1 Sepolia. Tx: %s", tx.Hash().Hex()))
	}(addr)

	// AUDIT
	go LogAction(addr, "ELECTION_ENDED", "Admin", "Manually ended election via API")

	respondJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Election ended successfully. Results are being anchored to L1."})
}

// GetAllElections returns a list of all elections (for Admin Dashboard)
func GetAllElections(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	if metadataCollection == nil {
		respondError(w, http.StatusInternalServerError, "DB not ready")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all, sort by StartDate desc
	opts := options.Find().SetSort(bson.M{"start_date": -1})
	cursor, err := metadataCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching elections")
		return
	}
	defer cursor.Close(ctx)

	var elections []ElectionMetadata
	if err = cursor.All(ctx, &elections); err != nil {
		respondError(w, http.StatusInternalServerError, "Error decoding elections")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   elections,
		"count":  len(elections),
	})
}

// GetArchivedResults fetches all anchored election results directly from the L1 Sepolia contract
func GetArchivedResults(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)

	l1Url := strings.TrimSpace(os.Getenv("L1_NODE_URL"))
	l1ArchiveAddr := strings.TrimSpace(os.Getenv("L1_ARCHIVE_CONTRACT_ADDRESS"))

	if l1Url == "" || l1ArchiveAddr == "" {
		respondError(w, http.StatusInternalServerError, "L1 Archiving is not configured")
		return
	}

	// 1. Connect to L1
	l1Client, err := ethclient.Dial(l1Url)
	if err != nil {
		log.Printf("[ARCHIVE ERR] Failed to connect to L1: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to connect to L1 Anchor")
		return
	}
	defer l1Client.Close()

	// 2. Bind Contract
	l1Archive, err := bindings.NewBindings(common.HexToAddress(l1ArchiveAddr), l1Client)
	if err != nil {
		log.Printf("[ARCHIVE ERR] Failed to bind L1 contract: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to bind L1 contract")
		return
	}

	// 3. To fetch all, we need the list of addresses.
	// Since mapping(address => FinalResult) isn't iterable in Solidity,
	// we will fetch all ENDED elections from MongoDB, and then query the L1 contract
	// for each address to get the verified on-chain result.

	if metadataCollection == nil {
		respondError(w, http.StatusInternalServerError, "DB not ready")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all ENDED elections
	cursor, err := metadataCollection.Find(ctx, bson.M{"status": "ENDED"})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching past elections")
		return
	}
	defer cursor.Close(ctx)

	var endedElections []ElectionMetadata
	if err = cursor.All(ctx, &endedElections); err != nil {
		respondError(w, http.StatusInternalServerError, "Error decoding past elections")
		return
	}

	// 4. Query L1 for each ended election
	type L1Result struct {
		ElectionAddress string `json:"election_address"`
		Title           string `json:"title"`
		WinnerName      string `json:"winner_name"`
		WinningVotes    int64  `json:"winning_votes"`
		TotalVoters     int64  `json:"total_voters"`
		Timestamp       int64  `json:"anchored_timestamp"`
	}

	var verifiedResults []L1Result
	callOpts := &bind.CallOpts{Context: context.Background(), Pending: false}

	for _, meta := range endedElections {
		addr := common.HexToAddress(meta.ElectionAddress)

		// Call archivedResults(address) mapping
		archived, err := l1Archive.ArchivedResults(callOpts, addr)
		if err != nil {
			log.Printf("Warning: Failed to fetch L1 result for %s: %v", meta.ElectionAddress, err)
			continue
		}

		// If the timestamp is 0, it means it hasn't been archived yet (or the archiving tx is still pending on L1)
		if archived.Timestamp == nil || archived.Timestamp.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		verifiedResults = append(verifiedResults, L1Result{
			ElectionAddress: archived.ElectionAddress.Hex(),
			Title:           archived.Title,
			WinnerName:      archived.WinnerName,
			WinningVotes:    archived.WinningVotes.Int64(),
			TotalVoters:     archived.TotalVoters.Int64(),
			Timestamp:       archived.Timestamp.Int64(),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   verifiedResults,
		"count":  len(verifiedResults),
	})
}
