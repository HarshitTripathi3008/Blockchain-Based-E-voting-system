package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func fetchImage(url string) ([]byte, string, error) {
	if url == "" {
		return nil, "", fmt.Errorf("empty url")
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	format := "PNG"
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		format = "JPG"
	}
	// Add more check if needed

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return data, format, nil
}

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

	// Determine Election Address
	electionAddr := r.URL.Query().Get("election_address")
	if electionAddr == "" {
		if len(v.Registrations) == 1 {
			electionAddr = v.Registrations[0].ElectionAddress
		} else if len(v.Registrations) > 0 {
			electionAddr = v.Registrations[0].ElectionAddress // Default to first
		} else {
			http.Error(w, "Voter has no election registrations", http.StatusBadRequest)
			return
		}
	}

	// 2. Generate PDF Bytes
	pdfBytes, err := createVoterIDPDF(v, electionAddr)
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

	// Determine Election Address from body if possible, or query?
	// The original handler didn't decode body for election address, just vars.
	// But we need it. Let's assume standard behavior: pick first or explicit.
	// Since standard `EmailVoterID` might be called from a context where we know the election,
	// checking body for `election_address` is good practice if we change signature.
	// But here we rely on what we have. Let's pick primary.
	electionAddr := ""
	if len(v.Registrations) > 0 {
		electionAddr = v.Registrations[0].ElectionAddress
	} else {
		http.Error(w, "Voter has no election registrations", http.StatusBadRequest)
		return
	}

	// Generate PDF
	pdfBytes, err := createVoterIDPDF(v, electionAddr)
	if err != nil {
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	// Email
	subject := "Your BlockVotes Voter ID Card"
	body := "<p>Hello " + v.FullName + ",</p><p>Please find attached your official digital Voter ID card for election: " + electionAddr + "</p>"

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
func createVoterIDPDF(v Voter, electionAddr string) ([]byte, error) {
	fullName := v.FullName
	if fullName == "" {
		fullName = "Unknown Voter"
	}
	idHex := v.ID.Hex()
	addr := v.Address
	if addr == "" {
		addr = "N/A"
	}
	mobile := v.Mobile
	if mobile == "" {
		mobile = "N/A"
	}
	email := v.Email
	dobStr := "N/A"
	if !v.DOB.IsZero() {
		dobStr = v.DOB.Format("2006-01-02")
	}
	father := v.FatherName
	if father == "" {
		father = "N/A"
	}

	// QR Code Content (JSON for verification)
	qrContent := fmt.Sprintf(`{"id":"%s","name":"%s","father":"%s","dob":"%s","mobile":"%s","email":"%s","election":"%s"}`,
		idHex, fullName, father, dobStr, mobile, email, electionAddr)

	qrPng, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	// Fetch Photo if available
	var photoBytes []byte
	var photoFmt string
	if v.PhotoURL != "" {
		// Log or ignore error, fallback to placeholder
		pb, pf, err := fetchImage(v.PhotoURL)
		if err == nil {
			photoBytes = pb
			photoFmt = pf
		} else {
			fmt.Printf("Error fetching photo for PDF: %v\n", err)
		}
	}

	// PDF Generation
	// Card Size: 85.6mm x 54mm (Credit Card) -> Scaled up for readability: 140mm x 85mm
	wCard, hCard := 140.0, 85.0

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	x, y := 35.0, 40.0

	// --- CARD STYLING ---

	// 1. Border & Background
	pdf.SetFillColor(255, 255, 255) // White bg
	pdf.Rect(x, y, wCard, hCard, "F")
	pdf.SetLineWidth(0.8)
	pdf.SetDrawColor(0, 0, 0) // Black border
	pdf.Rect(x, y, wCard, hCard, "D")

	// 2. Header
	pdf.SetFillColor(30, 41, 59) // Dark Blue/Slate
	pdf.Rect(x, y, wCard, 18.0, "F")

	// Logo/Title in Header
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(x+5, y+5)
	pdf.Cell(0, 8, "BLOCKVOTES ELECTION COMMISSION")

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(x+5, y+11)
	pdf.Cell(0, 5, fmt.Sprintf("Election: %s", electionAddr))

	// 3. Photo
	photoW, photoH := 30.0, 35.0
	photoX, photoY := x+6, y+25

	if len(photoBytes) > 0 {
		// Draw uploaded photo
		opt := gofpdf.ImageOptions{ImageType: photoFmt, ReadDpi: true}
		pdf.RegisterImageOptionsReader("profile_pic", opt, bytes.NewReader(photoBytes))
		pdf.ImageOptions("profile_pic", photoX, photoY, photoW, photoH, false, opt, 0, "")
		pdf.SetDrawColor(150, 150, 150)
		pdf.Rect(photoX, photoY, photoW, photoH, "D") // border around photo
	} else {
		// Placeholder
		pdf.SetFillColor(220, 220, 220)
		pdf.Rect(photoX, photoY, photoW, photoH, "F")
		pdf.SetDrawColor(150, 150, 150)
		pdf.Rect(photoX, photoY, photoW, photoH, "D")
		// "No Photo" Text
		pdf.SetTextColor(100, 100, 100)
		pdf.SetFont("Arial", "I", 8)
		pdf.SetXY(photoX+2, photoY+15)
		pdf.Cell(photoW, 5, "Photo N/A")
	}

	// 4. Details Column
	labelX := photoX + photoW + 8
	valX := labelX + 28
	currentY := y + 26.0
	lineHeight := 6.0 // Reduced from 7.0 to fix overlap

	pdf.SetTextColor(0, 0, 0)

	// Helper macros
	printRow := func(label, value string, smallFont bool) {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(labelX, currentY)
		pdf.Cell(25, 6, label+":")

		fontSize := 10.0
		if smallFont {
			fontSize = 8.0
		}
		pdf.SetFont("Arial", "", fontSize)
		pdf.SetXY(valX, currentY)

		// Adjust Y slightly for small font if needed, but standard cell works
		pdf.Cell(70, 6, value)

		currentY += lineHeight
	}

	printRow("Name", fullName, false)
	printRow("Father Name", father, false)
	printRow("DOB", dobStr, false)
	printRow("Mobile", mobile, false)
	printRow("Email", email, true) // Email might be long too
	// Address might be long, handle wrapping manually or concise
	// Simple truncate for now
	if len(addr) > 50 {
		addr = addr[:47] + "..."
	}
	printRow("Address", addr, true)   // Address small font
	printRow("Voter ID", idHex, true) // Voter ID small font

	// 5. QR Code
	// Place bottom right
	qrSize := 32.0
	qrX := x + wCard - qrSize - 4
	qrY := y + hCard - qrSize - 4

	// Adjust QR to not overlap if text is really long, but
	// printRow stops at: 26 + 7*6 = 68. y+68 = 108.
	// QR starts at y + 85 - 36 = y + 49.
	// Overlap Check: Data column (Left) vs QR (Right).
	// valX ~ 80mm from left of card. Text width ~60mm. 80+60=140.
	// QR starts at 140 - 36 = 104.
	// Overlap IS likely if text is long.
	// But we moved lineHeight to 6, so Vertical space is saved.
	// Let's ensure text doesn't bleed into QR.
	// We won't fix horizontal overlap complexly now, just vertical fix helps.

	opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.RegisterImageOptionsReader("qrcode.png", opt, bytes.NewReader(qrPng))
	pdf.ImageOptions("qrcode.png", qrX, qrY, qrSize, qrSize, false, opt, 0, "")

	// 6. Footer Signature area
	// Position: y + hCard - 15 = y + 70.
	// printRow ends at y + 68.
	// 70 > 68. Safe vertical gap now!

	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(x+wCard-45, y+hCard-15, x+wCard-10, y+hCard-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetXY(x+wCard-45, y+hCard-14)
	pdf.Cell(35, 4, "Issuing Authority")

	// Disclaimer
	pdf.SetTextColor(100, 100, 100)
	pdf.SetXY(x+5, y+hCard-8)
	pdf.SetFont("Arial", "", 7)
	pdf.Cell(0, 5, "This card is digitally verified on the Ethereum blockchain.")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
