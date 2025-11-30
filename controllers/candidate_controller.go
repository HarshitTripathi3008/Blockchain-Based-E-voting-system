package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"MAJOR-PROJECT/bindings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// CandidateRequest and CandidateDocument kept same as before (ensure candidateCollection init exists)
var candidateCollection *mongo.Collection

func InitCandidateCollection(client *mongo.Client, dbName string) {
	candidateCollection = client.Database(dbName).Collection("candidates")
	fmt.Println("✅ Initialized candidates collection")
}
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type CandidateRequest struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageHash       string `json:"imageHash"`
	ElectionName    string `json:"election_name,omitempty"`
	ElectionAddress string `json:"election_address,omitempty"`
}
type CandidateDocument struct {
	Name            string    `bson:"name"`
	Email           string    `bson:"email"`
	Description     string    `bson:"description"`
	ImageHash       string    `bson:"imageHash"`
	ElectionName    string    `bson:"electionName,omitempty"`
	ElectionAddress string    `bson:"electionAddress,omitempty"`
	TxHash          string    `bson:"txHash,omitempty"`
	Status          string    `bson:"status"` // e.g., "submitted", "mined", "reverted"
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}
// RegisterCandidate registers a candidate on-chain (using bindings.NewElection) and persists metadata.
func RegisterCandidate(w http.ResponseWriter, r *http.Request) {
	// CORS & headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CandidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "Invalid JSON body"})
		return
	}
	if req.Email == "" || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "email and name are required"})
		return
	}
	if req.ElectionAddress == "" {
		if addr := r.URL.Query().Get("election_address"); addr != "" {
			req.ElectionAddress = addr
		}
	}
	if req.ElectionAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "election_address is required"})
		return
	}

	client, err := getClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "Failed to connect to Ethereum node: " + err.Error()})
		return
	}
	defer client.Close()

	auth, err := getAuth()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "Failed to create transaction signer: " + err.Error()})
		return
	}

	contractAddr := common.HexToAddress(req.ElectionAddress)
	contract, err := bindings.NewElection(contractAddr, client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "Failed to bind to election contract: " + err.Error()})
		return
	}

	imgHash := req.ImageHash
	if imgHash == "" {
		imgHash = ""
	}

	tx, err := contract.AddCandidate(auth, req.Name, req.Description, imgHash, req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "Failed to register candidate on blockchain: " + err.Error()})
		return
	}

	txHashHex := tx.Hash().Hex()
	now := time.Now().UTC()

	// Persist candidate doc to MongoDB if candidateCollection initialized (same as before)
	if candidateCollection != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		doc := CandidateDocument{
			Name:            req.Name,
			Email:           req.Email,
			Description:     req.Description,
			ImageHash:       req.ImageHash,
			ElectionName:    req.ElectionName,
			ElectionAddress: req.ElectionAddress,
			TxHash:          txHashHex,
			Status:          "submitted",
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if _, err := candidateCollection.InsertOne(ctx, doc); err != nil {
			fmt.Printf("warning: failed to insert candidate into DB: %v\n", err)
		}
	}

	// optional wait for mining
	txTimeoutSec := 0
	if s := os.Getenv("TX_TIMEOUT"); s != "" {
		if t, err := strconv.Atoi(s); err == nil && t > 0 {
			txTimeoutSec = t
		}
	}
	if txTimeoutSec > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(txTimeoutSec)*time.Second)
		defer cancel()
		receipt, werr := bind.WaitMined(ctx, client, tx)
		if werr != nil {
			updateCandidateStatus(txHashHex, "pending")
			_ = json.NewEncoder(w).Encode(Response{Status: "pending", Message: "Transaction submitted but mining confirmation failed: " + werr.Error(), Data: map[string]interface{}{"txHash": txHashHex}})
		} else {
			if receipt == nil || receipt.Status != 1 {
				updateCandidateStatus(txHashHex, "reverted")
				_ = json.NewEncoder(w).Encode(Response{Status: "error", Message: "Transaction mined but reverted on-chain", Data: map[string]interface{}{"txHash": txHashHex}})
			} else {
				updateCandidateStatus(txHashHex, "mined")
				_ = json.NewEncoder(w).Encode(Response{Status: "success", Message: "Candidate registered and transaction mined", Data: map[string]interface{}{"txHash": txHashHex, "blockNumber": receipt.BlockNumber.String()}})
			}
		}
	} else {
		_ = json.NewEncoder(w).Encode(Response{Status: "success", Message: "Candidate registration transaction submitted", Data: map[string]interface{}{"txHash": txHashHex}})
	}

	_ = sendRegistrationEmail(req.Email, req.ElectionName)
}
func updateCandidateStatus(txHash, status string) {
	if candidateCollection == nil {
		fmt.Println("warning: candidateCollection is nil; cannot update status")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	filter := bson.M{"txHash": txHash}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().UTC()}}
	_, err := candidateCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("warning: failed to update candidate status for tx %s: %v\n", txHash, err)
	}
}

// sendRegistrationEmail sends a registration confirmation email using SMTP. It uses EMAIL and PASSWORD env vars.
func sendRegistrationEmail(to, electionName string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD") // keep env var name 'PASSWORD' as in your .env

	if from == "" || password == "" {
		return fmt.Errorf("email credentials not configured in environment")
	}

	const smtpHost = "smtp.gmail.com"
	const smtpPort = "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := fmt.Sprintf("%s Registration", electionName)
	body := fmt.Sprintf("Congratulations! You have been registered for the %s election.\n\nBest regards,\nVoting System Team", electionName)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
