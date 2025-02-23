package uploaders

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsBucket struct {
	S3Client *s3.Client
	Bucket   string `json:"bucket"`
	File     string `json:"file"`
	Key      string `json:"key"`
}

type UploadBackup interface {
	UploadBackup(context.Context, string) error
}

func InitAwsBucket(c context.Context, bucket, key string) (*AwsBucket, error) {
	config, err := config.LoadDefaultConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed loading config, %v", err)
	}
	s3Client := s3.NewFromConfig(config)
	return &AwsBucket{
		S3Client: s3Client,
		Bucket:   bucket,
		Key:      key,
	}, nil
}

func (b *AwsBucket) UploadBackup(c context.Context, filePath string) error {
	f, err := os.Open("/tmp/" + filePath)
	if err != nil {
		return fmt.Errorf("failed to open file, %v", err)
	}
	k := b.Key + filePath
	// Upload input parameters
	uploadParams := &s3.PutObjectInput{
		Bucket: &b.Bucket,
		Key:    &k,
		Body:   f,
	}

	ctx, cancel := context.WithTimeout(c, time.Second*50)
	defer cancel()
	result, err := b.S3Client.PutObject(ctx, uploadParams)
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	log.Printf("backup %s uploaded to, %v\n", b.Key, result)
	return nil
}
