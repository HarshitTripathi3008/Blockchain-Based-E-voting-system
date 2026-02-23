package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"MAJOR-PROJECT/util"
)

// Response returned to the frontend
type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	URL     string `json:"url,omitempty"`
}

// UnifiedUploadHandler implements the hybrid storage logic
func UnifiedUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Safety recovery
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("PANIC in UnifiedUploadHandler: %v", rec)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 10MB limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{Status: "error", Message: "invalid multipart form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{Status: "error", Message: "file is required"})
		return
	}
	defer file.Close()

	// Detect Content-Type from header (usually reliable from browser)
	mimeType := header.Header.Get("Content-Type")
	log.Printf("Processing upload: %s (Type: %s)", header.Filename, mimeType)

	var url string
	var uploadErr error

	// Logic: Images/Videos/Audio -> AWS S3, Else -> Google Drive
	isMedia := false
	if len(mimeType) >= 6 && (mimeType[:6] == "image/" || mimeType[:6] == "video/" || mimeType[:6] == "audio/") {
		isMedia = true
	}

	if isMedia {
		// Delegate to S3
		log.Println("Delegating to AWS S3...")
		url, uploadErr = util.UploadToS3(file, header.Filename)
	} else {
		// Delegate to Google Drive
		log.Println("Delegating to Google Drive...")
		url, uploadErr = util.UploadToGDrive(file, header.Filename)
	}

	if uploadErr != nil {
		log.Printf("Upload Error: %v", uploadErr)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UploadResponse{Status: "error", Message: "upload failed: " + uploadErr.Error()})
		return
	}

	log.Printf("Upload success: %s", url)

	// Return success
	json.NewEncoder(w).Encode(UploadResponse{
		Status:  "success",
		Message: "uploaded successfully",
		URL:     url,
	})
}
