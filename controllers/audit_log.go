package controllers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jung-kurt/gofpdf"
)

// AuditLog represents a single event in the election lifecycle
type AuditLog struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ElectionAddress string             `bson:"election_address" json:"election_address"`
	Action          string             `bson:"action" json:"action"`   // e.g. "ELECTION_CREATED", "VOTE_CAST", "ELECTION_ENDED"
	Actor           string             `bson:"actor" json:"actor"`     // email or ID of who performed it
	Details         string             `bson:"details" json:"details"` // extra info
	Timestamp       time.Time          `bson:"timestamp" json:"timestamp"`
}

var auditCollection *mongo.Collection

// InitAuditCollection initializes the audit_logs collection
func InitAuditCollection(client *mongo.Client, dbName string) {
	auditCollection = client.Database(dbName).Collection("audit_logs")
	fmt.Println("✅ Initialized audit_logs collection")
}

// LogAction records an event to the database
func LogAction(electionAddr, action, actor, details string) {
	if auditCollection == nil {
		log.Println("⚠️ Audit collection not initialized, skipping log:", action)
		return
	}

	entry := AuditLog{
		ElectionAddress: electionAddr,
		Action:          action,
		Actor:           actor,
		Details:         details,
		Timestamp:       time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := auditCollection.InsertOne(ctx, entry)
	if err != nil {
		log.Printf("❌ Failed to log action %s: %v", action, err)
	} else {
		log.Printf("📝 Audit Log: [%s] %s by %s", action, details, actor)
	}
}

// GetElectionLogs retrieves all logs for a specific election, sorted by time
func GetElectionLogs(electionAddr string) ([]AuditLog, error) {
	if auditCollection == nil {
		return nil, fmt.Errorf("audit collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}) // Oldest first for timeline
	cursor, err := auditCollection.Find(ctx, bson.M{"election_address": electionAddr}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GenerateAuditLogPDF creates a PDF of the audit logs
func GenerateAuditLogPDF(logs []AuditLog, electionName string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Election Audit Report: "+electionName)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, "Generated at: "+time.Now().Format(time.RFC1123))
	pdf.Ln(10)

	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(40, 7, "Time", "1", 0, "", true, 0, "")
	pdf.CellFormat(50, 7, "Action", "1", 0, "", true, 0, "")
	pdf.CellFormat(50, 7, "Actor", "1", 0, "", true, 0, "")
	pdf.CellFormat(50, 7, "Details", "1", 0, "", true, 0, "")
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 9)
	for _, l := range logs {
		pdf.CellFormat(40, 6, l.Timestamp.Format("15:04:05 02-Jan"), "1", 0, "", false, 0, "")
		pdf.CellFormat(50, 6, l.Action, "1", 0, "", false, 0, "")

		// Truncate actor if too long
		actor := l.Actor
		if len(actor) > 25 {
			actor = actor[:22] + "..."
		}
		pdf.CellFormat(50, 6, actor, "1", 0, "", false, 0, "")

		details := l.Details
		if len(details) > 25 {
			details = details[:22] + "..."
		}
		pdf.CellFormat(50, 6, details, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
