package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// LighthouseResponse loosely models common fields returned by Lighthouse / IPFS uploads
type LighthouseResponse struct {
	Cid  string `json:"cid,omitempty"`
	CID  string `json:"CID,omitempty"`
	Hash string `json:"Hash,omitempty"`
	// some APIs return Name/Size etc; we just care about CID/Hash
	Name string `json:"Name,omitempty"`
	Size int64  `json:"Size,omitempty"`
}

// UploadFileToLighthouse uploads a file to Lighthouse Storage and returns the CID/Hash.
func UploadFileToLighthouse(filePath string) (string, error) {
	apiKey := os.Getenv("LIGHTHOUSE_API_KEY")
	uploadURL := os.Getenv("LIGHTHOUSE_UPLOAD_URL")
	if uploadURL == "" {
		uploadURL = "https://node.lighthouse.storage/api/v0/add"
	}
	if apiKey == "" {
		log.Println("upload util: ERROR: LIGHTHOUSE_API_KEY is empty")
		return "", fmt.Errorf("missing LIGHTHOUSE_API_KEY")
	}
	log.Printf("upload util: uploading file %s to %s", filePath, uploadURL)
		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("cannot open file: %v", err)
		}
		defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("cannot create form file: %v", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("cannot copy file to form: %v", err)
	}
	// you can add optional form fields if Lighthouse supports them
	_ = writer.Close()

	req, err := http.NewRequest("POST", uploadURL, &body)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %v", err)
	}
	// Lighthouse expects Authorization: Bearer <key>
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("lighthouse request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("lighthouse returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var lresp LighthouseResponse
	if err := json.Unmarshal(respBody, &lresp); err != nil {
		// if unmarshalling fails, return raw body for debugging
		return "", fmt.Errorf("failed to parse lighthouse response: %v; body: %s", err, string(respBody))
	}

	// pick first available identifier
	cid := lresp.Cid
	if cid == "" {
		cid = lresp.CID
	}
	if cid == "" {
		cid = lresp.Hash
	}
	if cid == "" {
		// in some cases Lighthouse returns plain text or different JSON — return full response for debugging
		return "", fmt.Errorf("lighthouse response missing CID/Hash: %s", string(respBody))
	}

	return cid, nil
}
