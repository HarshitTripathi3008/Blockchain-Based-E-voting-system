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

// UploadToCloudinary handles file upload to Cloudinary
func UploadToCloudinary(w http.ResponseWriter, r *http.Request) {
	// Safety recovery
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("PANIC in UploadToCloudinary: %v", rec)
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

	log.Printf("Starting upload for file: %s", header.Filename)

	// Upload using the utility
	url, err := util.UploadToCloudinary(file, header.Filename)
	if err != nil {
		log.Printf("Cloudinary Upload Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UploadResponse{Status: "error", Message: "upload failed: " + err.Error()})
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
