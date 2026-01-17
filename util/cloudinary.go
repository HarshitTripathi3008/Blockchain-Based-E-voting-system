package util

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadToCloudinary uploads a file to Cloudinary and returns the secure URL.
func UploadToCloudinary(file multipart.File, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	// Upload file (strip extension from PublicID to avoid .png.png)
	publicID := filename
	if ext := filepath.Ext(filename); ext != "" {
		publicID = filename[:len(filename)-len(ext)]
	}

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "voting_system",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
