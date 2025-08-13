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

func (c *MinioClient) bucketExists(ctx context.Context, bucketname string) (bool, error) {
	// Check local cache for bucket
	if !slices.Contains(c.buckets, bucketname) {
		// Check if the bucket exists
		exists, err := c.client.BucketExists(ctx, bucketname)
		if err != nil {
			return false, fmt.Errorf("error checking bucket existence: %w", err)
		}

		if exists {
			// Add the bucket to the list of buckets
			c.buckets = append(c.buckets, bucketname)
		} else {
			return false, fmt.Errorf("bucket %s does not exist", bucketname)
		}
	}

	return true, nil
}

func (c *MinioClient) WriteBytesToFile(ctx context.Context, bucketname, objectName string, data []byte) error {
	// Check that a bucket name is in the list of buckets
	_, err := c.bucketExists(ctx, bucketname)
	if err != nil {
		return err
	}

	// Upload a file to the specified bucket
	_, err = c.client.PutObject(ctx, bucketname, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *MinioClient) AppendBytesToFile(ctx context.Context, bucketname, objectName string, data []byte) error {
	// Check that a bucket name is in the list of buckets
	_, err := c.bucketExists(ctx, bucketname)
	if err != nil {
		return err
	}

	// Append data to an existing file
	_, err = c.client.AppendObject(ctx, bucketname, objectName, bytes.NewReader(data), int64(len(data)), minio.AppendObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *MinioClient) ObjectExists(ctx context.Context, bucketname, objectName string) (bool, error) {
	_, err := c.client.StatObject(ctx, bucketname, objectName, minio.StatObjectOptions{})
	if err != nil {
		// Check if it's a "not found" error
		errorResponse := minio.ToErrorResponse(err)
		if errorResponse.Code == "NoSuchKey" {
			return false, nil // Object doesn't exist, but no error
		}
		return false, err // Some other error occurred (network, auth, etc.)
	}

	return true, nil // Object exists
}
