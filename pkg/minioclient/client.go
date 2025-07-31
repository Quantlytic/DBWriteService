package minioclient

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client  *minio.Client
	buckets []string
}

// NewMinioClient initializes a new Minio client with the provided configuration
func NewMinioClient(url, accessKeyID, secretAccessKey string) (*MinioClient, error) {
	// Create a new Minio client
	client, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, // Set to true if using HTTPS
	})

	if err != nil {
		return nil, err
	}

	return &MinioClient{
		client: client,
	}, nil
}

func (c *MinioClient) WriteBytesToFile(ctx context.Context, bucketname, objectName string, data []byte) error {
	// Check that a bucket name is in the list of buckets
	if !slices.Contains(c.buckets, bucketname) {
		// Check if the bucket exists
		exists, err := c.client.BucketExists(ctx, bucketname)
		if err != nil {
			fmt.Printf("Error checking bucket existence: %v\n", err)
			return err
		}

		if !exists {
			return fmt.Errorf("bucket %s does not exist", bucketname)
		}

		// Add the bucket to the list of buckets
		c.buckets = append(c.buckets, bucketname)
		fmt.Printf("Added bucket %s to the list of buckets\n", bucketname)
	}

	// Upload a file to the specified bucket
	_, err := c.client.PutObject(ctx, bucketname, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
