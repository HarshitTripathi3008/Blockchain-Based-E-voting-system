package controllers

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"MAJOR-PROJECT/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmailJobDocument struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	To             string             `bson:"to" json:"to"`
	Subject        string             `bson:"subject" json:"subject"`
	HTMLBody       string             `bson:"htmlBody" json:"htmlBody"`
	Filename       string             `bson:"filename,omitempty" json:"filename,omitempty"`
	AttachmentData []byte             `bson:"attachmentData,omitempty" json:"-"`
	Status         string             `bson:"status" json:"status"`
	Attempts       int                `bson:"attempts" json:"attempts"`
	LastError      string             `bson:"lastError,omitempty" json:"lastError,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	SentAt         *time.Time         `bson:"sentAt,omitempty" json:"sentAt,omitempty"`
}

type VoteJobDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ElectionAddress string             `bson:"electionAddress" json:"electionAddress"`
	CandidateID     int64              `bson:"candidateId" json:"candidateId"`
	VoterEmail      string             `bson:"voterEmail" json:"voterEmail"`
	Status          string             `bson:"status" json:"status"`
	TxHash          string             `bson:"txHash,omitempty" json:"txHash,omitempty"`
	LastError       string             `bson:"lastError,omitempty" json:"lastError,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	SubmittedAt     *time.Time         `bson:"submittedAt,omitempty" json:"submittedAt,omitempty"`
	MinedAt         *time.Time         `bson:"minedAt,omitempty" json:"minedAt,omitempty"`
}

var (
	emailJobCollection *mongo.Collection
	voteJobCollection  *mongo.Collection
	emailWakeCh        = make(chan struct{}, 1)
	voteWakeCh         = make(chan struct{}, 1)
)

func InitAsyncJobCollections(client *mongo.Client, dbName string) {
	db := client.Database(dbName)
	emailJobCollection = db.Collection("email_jobs")
	voteJobCollection = db.Collection("vote_jobs")

	go emailQueueWorker()
	go voteQueueWorker()
	log.Println("[OK] Initialized async job collections")
}

func notifyWorker(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func enqueueEmailJob(to, subject, htmlBody, filename string, attachmentData []byte) error {
	if emailJobCollection == nil {
		return fmt.Errorf("email job collection not initialized")
	}

	now := time.Now().UTC()
	doc := EmailJobDocument{
		To:             to,
		Subject:        subject,
		HTMLBody:       htmlBody,
		Filename:       filename,
		AttachmentData: attachmentData,
		Status:         "pending",
		Attempts:       0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := emailJobCollection.InsertOne(ctx, doc); err != nil {
		return err
	}

	notifyWorker(emailWakeCh)
	return nil
}

func queueEmailWithAttachment(to, subject, htmlBody, filename string, attachmentData []byte) error {
	return enqueueEmailJob(to, subject, htmlBody, filename, attachmentData)
}

func claimNextEmailJob(ctx context.Context) (*EmailJobDocument, error) {
	if emailJobCollection == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	opts := options.FindOneAndUpdate().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetReturnDocument(options.After)

	var job EmailJobDocument
	err := emailJobCollection.FindOneAndUpdate(
		ctx,
		bson.M{"status": "pending"},
		bson.M{"$set": bson.M{"status": "processing", "updatedAt": now}},
		opts,
	).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func emailQueueWorker() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-emailWakeCh:
		case <-ticker.C:
		}

		for {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			job, err := claimNextEmailJob(ctx)
			cancel()
			if err != nil {
				log.Printf("[EMAIL WORKER] claim error: %v", err)
				break
			}
			if job == nil {
				break
			}

			var sendErr error
			for attempt := 1; attempt <= 3; attempt++ {
				sendErr = sendEmailWithAttachment(job.To, job.Subject, job.HTMLBody, job.Filename, job.AttachmentData)
				if sendErr == nil {
					break
				}
				log.Printf("[EMAIL RETRY %d/3] %s: %v", attempt, job.To, sendErr)
				if attempt < 3 {
					time.Sleep(time.Duration(attempt) * 2 * time.Second)
				}
			}

			updateCtx, updateCancel := context.WithTimeout(context.Background(), 10*time.Second)
			now := time.Now().UTC()
			update := bson.M{
				"$set": bson.M{
					"updatedAt": now,
					"attempts":  job.Attempts + 1,
				},
			}
			if sendErr == nil {
				update["$set"].(bson.M)["status"] = "sent"
				update["$set"].(bson.M)["sentAt"] = now
				update["$set"].(bson.M)["lastError"] = ""
			} else {
				update["$set"].(bson.M)["status"] = "failed"
				update["$set"].(bson.M)["lastError"] = sendErr.Error()
			}
			if _, err := emailJobCollection.UpdateByID(updateCtx, job.ID, update); err != nil {
				log.Printf("[EMAIL WORKER] update error for %s: %v", job.ID.Hex(), err)
			}
			updateCancel()
		}
	}
}

func enqueueVoteJob(electionAddress string, candidateID int64, voterEmail string) (*VoteJobDocument, error) {
	if voteJobCollection == nil {
		return nil, fmt.Errorf("vote job collection not initialized")
	}

	now := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existing VoteJobDocument
	err := voteJobCollection.FindOne(ctx, bson.M{
		"electionAddress": electionAddress,
		"voterEmail":      voterEmail,
		"status": bson.M{
			"$in": []string{"queued", "processing", "submitted", "mined"},
		},
	}).Decode(&existing)
	if err == nil {
		return &existing, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	doc := VoteJobDocument{
		ElectionAddress: electionAddress,
		CandidateID:     candidateID,
		VoterEmail:      voterEmail,
		Status:          "queued",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	result, err := voteJobCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	doc.ID = result.InsertedID.(primitive.ObjectID)
	notifyWorker(voteWakeCh)
	return &doc, nil
}

func getVoteJobByVoter(electionAddr, voterEmail string) (*VoteJobDocument, error) {
	if voteJobCollection == nil {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job VoteJobDocument
	err := voteJobCollection.FindOne(ctx, bson.M{
		"electionAddress": electionAddr,
		"voterEmail":      voterEmail,
	}, options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &job, err
}

func claimNextVoteJob(ctx context.Context) (*VoteJobDocument, error) {
	if voteJobCollection == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	opts := options.FindOneAndUpdate().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetReturnDocument(options.After)

	var job VoteJobDocument
	err := voteJobCollection.FindOneAndUpdate(
		ctx,
		bson.M{"status": "queued"},
		bson.M{"$set": bson.M{"status": "processing", "updatedAt": now}},
		opts,
	).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func voteQueueWorker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-voteWakeCh:
		case <-ticker.C:
		}

		for {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			job, err := claimNextVoteJob(ctx)
			cancel()
			if err != nil {
				log.Printf("[VOTE WORKER] claim error: %v", err)
				break
			}
			if job == nil {
				break
			}

			if err := processVoteJob(job); err != nil {
				log.Printf("[VOTE WORKER] process error for %s: %v", job.ID.Hex(), err)
			}
		}
	}
}

func processVoteJob(job *VoteJobDocument) error {
	client, err := getClient()
	if err != nil {
		return updateVoteJobFailure(job.ID, err)
	}

	contractAddr := common.HexToAddress(job.ElectionAddress)
	if err := ensureContractVerified(client, contractAddr, "election"); err != nil {
		return updateVoteJobFailure(job.ID, err)
	}

	contract, err := bindings.NewElection(contractAddr, client)
	if err != nil {
		return updateVoteJobFailure(job.ID, err)
	}

	// --- CONCURRENCY CHECK: Is another job already active for this voter? ---
	// This prevents multiple submissions for the same voter before the first one is mined.
	count, err := voteJobCollection.CountDocuments(context.Background(), bson.M{
		"electionAddress": job.ElectionAddress,
		"voterEmail":      job.VoterEmail,
		"status":          bson.M{"$in": []string{"submitted", "processing", "pending_confirmation"}},
		"_id":             bson.M{"$ne": job.ID}, // exclude self
	})
	if err == nil && count > 0 {
		log.Printf("[SKIP] Another active job already exists for voter %s. Failing this duplicate.", job.VoterEmail)
		return updateVoteJobFailure(job.ID, fmt.Errorf("duplicate vote job detected"))
	}

	// --- SAFETY PRE-FLIGHT CHECK ---
	// Check if the voter has already voted via a call (free) before sending a transaction (paid).
	// This avoids "execution reverted: Error: You cannot double vote" on the blockchain.
	voterInfo, err := contract.Voters(&bind.CallOpts{Context: context.Background()}, job.VoterEmail)
	if err == nil && voterInfo.Voted {
		log.Printf("[WARN] Voter %s already voted for candidate %d. Failing job.", job.VoterEmail, voterInfo.CandidateIdVoted.Int64())
		return updateVoteJobFailure(job.ID, fmt.Errorf("voter has already voted"))
	}

	log.Printf("[INFO] Submitting vote to blockchain: Election=%s, Voter=%s, CandidateID=%d", job.ElectionAddress, job.VoterEmail, job.CandidateID)

	tx, err := submitChainTx(
		client,
		func() (*bind.TransactOpts, error) {
			return getAuth(client)
		},
		func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return contract.Vote(auth, big.NewInt(job.CandidateID), job.VoterEmail)
		},
	)
	if err != nil {
		log.Printf("[ERROR] Vote submission failed for %s: %v", job.VoterEmail, err)
		return updateVoteJobFailure(job.ID, err)
	}

	now := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := voteJobCollection.UpdateByID(ctx, job.ID, bson.M{
		"$set": bson.M{
			"status":      "submitted",
			"txHash":      tx.Hash().Hex(),
			"submittedAt": now,
			"updatedAt":   now,
			"lastError":   "",
		},
	}); err != nil {
		return err
	}

	go trackVoteMining(job.ID, job.ElectionAddress, job.VoterEmail, tx)

	// PERMANENT LOCK: Update voter status in DB as soon as transaction is SUBMITTED
	// this makes the "hasVoted" check true even before the tx is mined.
	if voterCollection != nil {
		vCtx, vCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer vCancel()
		_, _ = voterCollection.UpdateOne(vCtx,
			bson.M{"email": job.VoterEmail, "registrations.election_address": job.ElectionAddress},
			bson.M{"$set": bson.M{"registrations.$.status": "Voted"}},
		)
	}

	return nil
}

func updateVoteJobFailure(id primitive.ObjectID, err error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, updateErr := voteJobCollection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"status":    "failed",
			"lastError": err.Error(),
			"updatedAt": time.Now().UTC(),
		},
	})
	if updateErr != nil {
		return fmt.Errorf("%v; update failure: %w", err, updateErr)
	}
	return err
}

func trackVoteMining(id primitive.ObjectID, electionAddress, voterEmail string, tx *types.Transaction) {
	client, err := getClient()
	if err != nil {
		_ = updateVoteJobFailure(id, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	receipt, waitErr := bind.WaitMined(ctx, client, tx)
	now := time.Now().UTC()
	update := bson.M{"$set": bson.M{"updatedAt": now}}
	switch {
	case waitErr != nil:
		update["$set"].(bson.M)["status"] = "pending_confirmation"
		update["$set"].(bson.M)["lastError"] = waitErr.Error()
	case receipt == nil || receipt.Status != 1:
		update["$set"].(bson.M)["status"] = "reverted"
		update["$set"].(bson.M)["lastError"] = "vote transaction reverted"
	default:
		update["$set"].(bson.M)["status"] = "mined"
		update["$set"].(bson.M)["minedAt"] = now
		update["$set"].(bson.M)["lastError"] = ""

		// PERMANENT FIX: Update the voter's registration status to "Voted" in MongoDB
		if voterCollection != nil {
			vCtx, vCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer vCancel()
			_, _ = voterCollection.UpdateOne(vCtx,
				bson.M{"email": voterEmail, "registrations.election_address": electionAddress},
				bson.M{"$set": bson.M{"registrations.$.status": "Voted"}},
			)
		}

		invalidateCachePrefix(cacheKey("election", strings.ToLower(electionAddress)))
		go LogAction(electionAddress, "VOTE_CAST", voterEmail, "Voted successfully (mined)")
	}

	updateCtx, updateCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer updateCancel()
	if _, err := voteJobCollection.UpdateByID(updateCtx, id, update); err != nil {
		log.Printf("[VOTE WORKER] mining update error for %s: %v", id.Hex(), err)
	}
}

func GetVoteJobStatus(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)

	if voteJobCollection == nil {
		respondError(w, http.StatusInternalServerError, "vote job collection not initialized")
		return
	}

	jobID := mux.Vars(r)["jobId"]
	objID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid vote job id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var job VoteJobDocument
	if err := voteJobCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&job); err != nil {
		if err == mongo.ErrNoDocuments {
			respondError(w, http.StatusNotFound, "vote job not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to load vote job")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   job,
	})
}

// cache helpers
type cacheEntry struct {
	payload   interface{}
	expiresAt time.Time
}

var readCache sync.Map

func getCachedValue(key string) (interface{}, bool) {
	if v, ok := readCache.Load(key); ok {
		entry := v.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.payload, true
		}
		readCache.Delete(key)
	}
	return nil, false
}

func setCachedValue(key string, payload interface{}, ttl time.Duration) {
	readCache.Store(key, cacheEntry{payload: payload, expiresAt: time.Now().Add(ttl)})
}

func invalidateCachePrefix(prefix string) {
	readCache.Range(func(key, _ interface{}) bool {
		s, ok := key.(string)
		if ok && strings.HasPrefix(s, prefix) {
			readCache.Delete(key)
		}
		return true
	})
}

func cacheKey(parts ...string) string {
	return strings.Join(parts, "|")
}
