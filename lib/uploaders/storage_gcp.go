package uploaders

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// GcpBucket defines the interface for Google Cloud Storage bucket operations.
type GcpBucket struct {
	StorageClient *storage.Client
	Bucket        string `json:"bucket"`
	File          string `json:"file"`
	Key           string `json:"key"`
}

func InitGcpBucket(c context.Context, bucket, key string) (*GcpBucket, error) {
	client, err := storage.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()
	return &GcpBucket{
		StorageClient: client,
		Bucket:        bucket,
		Key:           key,
	}, nil
}

func (b *GcpBucket) UploadBackup(c context.Context, filePath string) error {

	// Open local file.
	f, err := os.Open("/tmp/" + filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(c, time.Second*50)
	defer cancel()

	client := *b.StorageClient
	bucket := b.Bucket
	object := b.Key + filePath
	bucketObj := client.Bucket(bucket).Object(object)

	// Upload an object with storage.Writer.
	wc := bucketObj.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	log.Printf("Blob %v uploaded.\n", object)
	return nil
}
