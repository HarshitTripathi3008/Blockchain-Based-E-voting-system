package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	// Status derived derived from time or manual override
	Status string `bson:"status" json:"status"` // "UPCOMING", "ONGOING", "ENDED"
}

var metadataCollection *mongo.Collection

// InitMetadataCollection initializes the collection
func InitMetadataCollection(client *mongo.Client, dbName string) {
	metadataCollection = client.Database(dbName).Collection("election_metadata")
	fmt.Println("✅ Initialized election_metadata collection")
}

// EnsureMetadata creates a default metadata entry if one doesn't exist, or extends it if expired.
func EnsureMetadata(electionAddr string) {
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
		}
		metadataCollection.InsertOne(ctx, newMeta)
		fmt.Printf("✅ Created metadata for %s (Expires: %s)\n", electionAddr, newMeta.EndDate)

	case nil:
		// Exists - Check if expired (Strict check, no auto-extend)
		if now.After(meta.EndDate) && meta.Status != "ENDED" {
			// Mark as ENDED if technically expired but not marked yet?
			// Or just leave it as ONGOING but IsElectionActive will return false.
			// Let's leave it, but ensure we DON'T extend.
			fmt.Printf("ℹ️ Election %s is past EndDate (%s)\n", electionAddr, meta.EndDate)
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

	// AUDIT
	go LogAction(addr, "ELECTION_ENDED", "Admin", "Manually ended election via API")

	respondJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Election ended successfully"})
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
