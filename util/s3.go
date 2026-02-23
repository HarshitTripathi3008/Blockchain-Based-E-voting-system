package util

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadToS3 uploads a file to S3 and returns the public URL.
func UploadToS3(file multipart.File, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")

	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
		return "", fmt.Errorf("S3 credentials OR bucket config missing in .env")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// Ensure unique filename
	objectKey := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), filename)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   file,
		// If the bucket is public-read, you can omit ACL or set it if needed
		// ACL: types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %v", err)
	}

	// Assuming the bucket has public read policy or is served via CloudFront
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, objectKey)
	return url, nil
}

// DeleteFromS3 deletes a file from S3 by its key.
func DeleteFromS3(objectKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")

	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
		return fmt.Errorf("S3 configuration missing")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	return err
}
