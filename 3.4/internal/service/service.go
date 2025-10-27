package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"threeFour/domain"
	"threeFour/internal/db"
	"threeFour/pkg"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/kafka"
)

type Serv struct {
	logger   zerolog.Logger
	producer *kafka.Producer
	consumer *kafka.Consumer
	storage  domain.MinioInt
	db       db.DB
}

func NewService(ctx context.Context, producer *kafka.Producer, consumer *kafka.Consumer, logger zerolog.Logger, storage domain.MinioInt, db db.DB) *Serv {
	return &Serv{
		logger:   logger,
		producer: producer,
		consumer: consumer,
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
	task := domain.Task{
		ID: id, Filename: image.Filename, Bucket: "uploads",
	}
	data, _ := json.Marshal(task)
	err = s.producer.Send(ctx, []byte("tasks"), data)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Serv) Get(ctx context.Context, id string) ([][]byte, error) {

	image, err := s.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if image.Status != "completed" {
		return nil, nil
	}
	if image.ResizedPath == nil {
		return nil, fmt.Errorf("processed path is null for image %s", id)
	}
	var result [][]byte
	find := []*string{image.ResizedPath, image.WatermarkPath, image.ThumbnailPath}
	for _, r := range find {
		data, err := s.storage.GetImage(ctx, "processed", *r)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, nil
}

func (s *Serv) Delete(ctx context.Context, id string) error {
	image, err := s.db.Get(ctx, id)
	if err != nil {
		return err
	}
	err = s.storage.Delete(ctx, image.OriginalPath)
	if err != nil {
		return err
	}

	err = s.db.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Serv) Update(ctx context.Context, id string, proccPath string) error {
	return nil
}

func (s *Serv) Convee(ctx context.Context) {
	for {
		msg, err := s.consumer.Fetch(ctx)
		if err != nil {
			s.logger.Err(err).Msg("Ошибка получения задачи")
			continue
		}

		fmt.Printf("Получена задача: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
		var task domain.Task
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			s.logger.Err(err).Msg("Ошибка разбора JSON задачи")
			continue
		}
		data, err := s.storage.GetImage(ctx, task.Bucket, task.Filename)
		if err != nil {
			s.logger.Err(err).Msg("Ошибка получения изображения")
			continue
		}

		reader := bytes.NewReader(data)
		img, format, err := image.Decode(reader)
		if err != nil {
			s.logger.Err(err).Msg("Ошибка декодирования изображения")
			continue
		}

		outputs, err := pkg.ProcessImage(img, "/assets/watermark.png", format)
		if err != nil {
			s.logger.Err(err).Msg("Ошибка обработки изображения")
			continue
		}
		paths := make([]string, 3)

		for i, r := range outputs {
			var name string
			switch i {
			case 0:
				name = "thumb_" + task.Filename
			case 1:
				name = "resized_" + task.Filename
			case 2:
				name = "watermarked_" + task.Filename
			}

			if buf, ok := r.(*bytes.Buffer); ok {
				size := buf.Len()
				_, err := s.storage.UploadToProcessed(ctx, name, buf, int64(size), format)
				if err != nil {
					s.logger.Err(err).Msg("Ошибка загрузки (buffer)")
					continue
				}
				paths[i] = name
			} else {
				data, err := io.ReadAll(r)
				if err != nil {
					s.logger.Err(err).Msg("Ошибка чтения из r")
					continue
				}
				size := len(data)
				path, err := s.storage.UploadToProcessed(ctx, name, bytes.NewReader(data), int64(size), format)
				if err != nil {
					s.logger.Err(err).Msg("Ошибка загрузки (reader)")
					continue
				}
				paths[i] = path
			}
		}

		if len(paths) == 3 {
			err = s.db.UpdatePathsAndStatus(ctx, task.ID, paths[0], paths[1], paths[2])
			if err != nil {
				s.logger.Err(err).Msg("Ошибка обновления статуса или конечного пути в DB")
			}
		}

		if err := s.consumer.Commit(ctx, msg); err != nil {
			s.logger.Err(err).Msg("Ошибка подтверждения задачи")
		}
		if err := s.db.Update(ctx, task.ID); err != nil {
			s.logger.Err(err).Msg("Ошибка обновления статуса задачи")
		}
		s.logger.Info().Msgf("Задача завершена: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
	}
}
