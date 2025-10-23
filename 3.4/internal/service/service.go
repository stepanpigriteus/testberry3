package service

import (
	"bytes"
	"context"
	"threeFour/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/kafka"
)

type Serv struct {
	logger   zerolog.Logger
	producer *kafka.Producer
	storage  domain.MinioInt
}

func NewService(ctx context.Context, producer *kafka.Producer, logger zerolog.Logger, storage domain.MinioInt) *Serv {
	return &Serv{
		logger:   logger,
		producer: producer,
		storage:  storage,
	}
}

func (s *Serv) Upload(ctx context.Context, image domain.ImageData) error {
	reader := bytes.NewReader(image.Bytes)
	fileSize := int64(len(image.Bytes))

	s.storage.UploadToUploads(ctx, image.Filename, reader, fileSize, image.ContentType)

	return nil
}

func (s *Serv) Get(ctx context.Context, id int) (domain.ImageData, error) {
	var image domain.ImageData
	return image, nil
}

func (s *Serv) Delete(ctx context.Context, id int) error {
	return nil
}
