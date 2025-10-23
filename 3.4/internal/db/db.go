package db

import (
	"context"
	"database/sql"
	"fmt"
	"threeFour/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type DB struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewDb(ctx context.Context, masterDSN string, slaveDSNs []string, logger zerolog.Logger) *DB {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("init database error %s", err)
	}
	return &DB{
		db:     db,
		logger: logger,
	}
}

func (d *DB) Create(ctx context.Context, original string) (string, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO images (id, original_path, status)
		VALUES ($1, $2, $3)
	`

	_, err := d.db.Master.ExecContext(ctx, query, id, original, "pending")
	if err != nil {
		return "", err
	}

	return id, nil
}

func (d *DB) Get(ctx context.Context, id string) (domain.ImageStat, error) {
	query := `
        SELECT id, original_path, processed_path, status, created_at
        FROM images
        WHERE id = $1
    `

	var image domain.ImageStat
	var processedPath sql.NullString

	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.OriginalPath,
		&processedPath,
		&image.Status,
		&image.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ImageStat{}, fmt.Errorf("image not found: %s", id)
		}
		return domain.ImageStat{}, fmt.Errorf("failed to get image: %w", err)
	}

	if processedPath.Valid {
		image.ProcessedPath = &processedPath.String
	} else {
		image.ProcessedPath = nil
	}

	return image, nil
}
