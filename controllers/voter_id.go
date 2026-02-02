package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateVoterID generates a PDF Voter ID card
// GET /api/voters/{voterId}/card
func GenerateVoterID(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	voterID := vars["voterId"]
	if voterID == "" {
		http.Error(w, "Voter ID required", http.StatusBadRequest)
		return
	}

	// 1. Fetch Voter
	if voterCollection == nil {
		http.Error(w, "DB not initialized", http.StatusInternalServerError)
		return
	}
	objID, _ := primitive.ObjectIDFromHex(voterID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var v Voter
	err := voterCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&v)
	if err != nil {
		http.Error(w, "Voter not found", http.StatusNotFound)
		return
	}

	// 2. Generate PDF Bytes
	pdfBytes, err := createVoterIDPDF(v)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// 3. Output
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=VoterID_%s.pdf", v.ID.Hex()))
	w.Write(pdfBytes)
}

// EmailVoterID generates PDF and emails it to the voter
func EmailVoterID(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	voterID := vars["voterId"]
	if voterID == "" {
		http.Error(w, "Voter ID required", http.StatusBadRequest)
		return
	}

	if voterCollection == nil {
		http.Error(w, "DB not initialized", http.StatusInternalServerError)
		return
	}
	objID, _ := primitive.ObjectIDFromHex(voterID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var v Voter
	err := voterCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&v)
	if err != nil {
		http.Error(w, "Voter not found", http.StatusNotFound)
		return
	}

	// Generate PDF
	pdfBytes, err := createVoterIDPDF(v)
	if err != nil {
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	// Email
	subject := "Your BlockVotes Voter ID Card"
	body := "<p>Hello " + v.FullName + ",</p><p>Please find attached your official digital Voter ID card.</p>"

	// Create attachment
	// NOTE: We need to import "github.com/sendgrid/sendgrid-go/helpers/mail" in this file
	// OR we can move sendEmail to a shared utils package.
	// For now, assuming sendEmail is available in 'controllers' package (same package)
	// AND sendEmail signature has been updated to accept attachments.

	// Since we are in the same package 'controllers', we can call sendEmail directly.
	// We will pass the attachment struct.

	// Construct attachment
	// We need to base64 encode the content?
	// The sendgrid helper handles it if we use mail.NewAttachment.
	// But since we can't easily import 'mail' here effectively if it wasn't already,
	// let's assume we update sendEmail to take raw bytes and filename, handling the mail struct construction inside sendEmail.

	// Wait, sendEmail is in voter.go. I will update it to take `attachments ...AttachmentData`.
	// Let's define a simple struct for passing data.

	err = sendEmailWithAttachment(v.Email, subject, body, "VoterID.pdf", pdfBytes)
	if err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"success", "message":"Email sent successfully"}`))
}

// Helper to create PDF bytes
func createVoterIDPDF(v Voter) ([]byte, error) {
	fullName := v.FullName
	if fullName == "" {
		fullName = "Unknown Voter"
	}
	idHex := v.ID.Hex()
	addr := v.Address
	if addr == "" {
		addr = "N/A"
	}
	dobStr := "N/A"
	if !v.DOB.IsZero() {
		dobStr = v.DOB.Format("2006-01-02")
	}

	// QR Code
	qrContent := fmt.Sprintf("BlockVotes Verified Voter\nID: %s\nName: %s\nElection: %s", idHex, fullName, v.ElectionAddress)
	qrPng, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	// PDF Generation
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	x, y := 40.0, 40.0
	wCard, hCard := 130.0, 80.0

	// Background & Border
	pdf.SetFillColor(245, 247, 250)
	pdf.Rect(x, y, wCard, hCard, "F")
	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(100, 100, 100)
	pdf.Rect(x, y, wCard, hCard, "D")

	// Header
	pdf.SetFillColor(108, 92, 231)
	pdf.Rect(x, y, wCard, 15.0, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(x+5, y+3)
	pdf.Cell(0, 10, "BLOCKVOTES VOTER ID")

	// Content
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(x+5, y+25)
	pdf.Cell(0, 8, fullName)

	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(x+5, y+35)
	pdf.Cell(0, 6, fmt.Sprintf("ID: %s", idHex))
	pdf.SetXY(x+5, y+41)
	pdf.Cell(0, 6, fmt.Sprintf("DOB: %s", dobStr))
	pdf.SetXY(x+5, y+47)
	pdf.Cell(0, 6, fmt.Sprintf("Loc: %s", addr))

	// Embed QR
	opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.RegisterImageOptionsReader("qrcode.png", opt, bytes.NewReader(qrPng))
	pdf.ImageOptions("qrcode.png", x+wCard-35, y+25, 30, 30, false, opt, 0, "")

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetXY(x+5, y+hCard-10)
	pdf.Cell(0, 5, "Official Digital Voter Card")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
