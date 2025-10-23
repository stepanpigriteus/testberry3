package service

import (
	"bytes"
	"context"
	"fmt"
	"threeFour/domain"
	"threeFour/internal/db"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/kafka"
)

type Serv struct {
	logger   zerolog.Logger
	producer *kafka.Producer
	storage  domain.MinioInt
	db       db.DB
}

func NewService(ctx context.Context, producer *kafka.Producer, logger zerolog.Logger, storage domain.MinioInt, db db.DB) *Serv {
	return &Serv{
		logger:   logger,
		producer: producer,
		storage:  storage,
		db:       db,
	}
}

func (s *Serv) Upload(ctx context.Context, image domain.ImageData) (string, error) {
	reader := bytes.NewReader(image.Bytes)
	fileSize := int64(len(image.Bytes))

	url, err := s.storage.UploadToUploads(ctx, image.Filename, reader, fileSize, image.ContentType)
	if err != nil {
		s.logger.Err(err).Msg("failed to upload in upStorage")
		return "", err
	}

	id, err := s.db.Create(ctx, url)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Serv) Get(ctx context.Context, id string) ([]byte, error) {
	image, err := s.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if image.ProcessedPath == nil {
		return nil, fmt.Errorf("processed path is null for image %s", id)
	}
	fmt.Println(image.ProcessedPath)
	data, err := s.storage.GetImage(ctx, "processed", *image.ProcessedPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Serv) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *Serv) Update(ctx context.Context, id string, proccPath string) error {
	return nil
}
