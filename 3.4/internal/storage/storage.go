package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client          *minio.Client
	UploadsBucket   string
	ProcessedBucket string
}

func NewMinioStorage(endpoint, accessKey, secretKey string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	storage := &MinioStorage{
		client:          client,
		UploadsBucket:   "uploads",
		ProcessedBucket: "processed",
	}

	ctx := context.Background()
	for _, bucket := range []string{storage.UploadsBucket, storage.ProcessedBucket} {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("check bucket %s: %w", bucket, err)
		}
		if !exists {
			if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return nil, fmt.Errorf("create bucket %s: %w", bucket, err)
			}
		}
	}

	return storage, nil
}

func (s *MinioStorage) UploadToUploads(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	return s.upload(ctx, s.UploadsBucket, objectName, file, fileSize, contentType)
}

func (s *MinioStorage) UploadToProcessed(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	return s.upload(ctx, s.ProcessedBucket, objectName, file, fileSize, contentType)
}

func (m *MinioStorage) GetImage(ctx context.Context, bucket, filename string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return data, nil
}

func (s *MinioStorage) upload(ctx context.Context, bucket, object string, file io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, bucket, object, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload to bucket %s failed: %w", bucket, err)
	}

	url := object
	return url, nil
}

func (s *MinioStorage) Delete(ctx context.Context, filename string) error {
	if err := s.deleteFromBucket(ctx, s.UploadsBucket, filename); err != nil {
		return fmt.Errorf("delete from uploads: %w", err)
	}

	suffixes := []string{"Res", "Min", "Water"}
	for _, suffix := range suffixes {
		name := suffix + filename
		if err := s.deleteFromBucket(ctx, s.ProcessedBucket, name); err != nil {
			return fmt.Errorf("delete from processed (%s): %w", name, err)
		}
	}

	return nil
}

func (s *MinioStorage) deleteFromBucket(ctx context.Context, bucket, object string) error {
	err := s.client.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object %s from bucket %s: %w", object, bucket, err)
	}
	return nil
}
