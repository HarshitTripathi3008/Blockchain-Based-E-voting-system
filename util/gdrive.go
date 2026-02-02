package util

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// UploadToGDrive uploads a file to Google Drive and returns the web view link.
// It requires a service account JSON key file at "gdrive_credentials.json"
// or passed via GDRIVE_CREDENTIALS_JSON environment variable.
func UploadToGDrive(file multipart.File, filename string) (string, error) {
	ctx := context.Background()

	// 1. Get Credentials
	creds := os.Getenv("GDRIVE_CREDENTIALS_JSON")
	var opt option.ClientOption
	if creds != "" {
		opt = option.WithCredentialsJSON([]byte(creds))
	} else {
		// Fallback to file
		opt = option.WithCredentialsFile("gdrive_credentials.json")
	}

	// 2. Create Service
	srv, err := drive.NewService(ctx, opt)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	// 3. Create File Metadata
	f := &drive.File{
		Name: filename,
		// Parents: []string{"folder_id"}, // Optional: if we want to put it in a specific folder
	}

	// 4. Create and Upload
	res, err := srv.Files.Create(f).Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}

	// 5. Make Public (Reader)
	perm := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err = srv.Permissions.Create(res.Id, perm).Do()
	if err != nil {
		return "", fmt.Errorf("unable to permission file: %v", err)
	}

	// 6. Get WebViewLink
	// We need to fetch the file again to get the link, or specify fields in Create
	// But Create returns *File, let's check if it has WebViewLink populated.
	// Usually it doesn't unless requested.
	// Let's fetch it explicitly.
	fileInfo, err := srv.Files.Get(res.Id).Fields("webViewLink").Do()
	if err != nil {
		return "", fmt.Errorf("unable to get file info: %v", err)
	}

	return fileInfo.WebViewLink, nil
}
