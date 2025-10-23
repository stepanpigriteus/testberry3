package domain

import (
	"context"
	"io"
)

type MinioInt interface {
	UploadToUploads(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
	UploadToProcessed(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
}
