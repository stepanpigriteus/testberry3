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
	fmt.Println(id)
	query := `
        SELECT id, original_path,  thumbnail_path,watermark_path ,resized_path,  status, created_at
        FROM images
        WHERE id = $1
    `

	var image domain.ImageStat

	var thumb, watermark, resized sql.NullString

	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.OriginalPath,
		&thumb,
		&watermark,
		&resized,
		&image.Status,
		&image.CreatedAt,
	)

	if thumb.Valid {
		image.ThumbnailPath = &thumb.String
	}
	if watermark.Valid {
		image.WatermarkPath = &watermark.String
	}
	if resized.Valid {
		image.ResizedPath = &resized.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ImageStat{}, fmt.Errorf("image not found: %s", id)
		}
		return domain.ImageStat{}, fmt.Errorf("failed to get image: %w", err)
	}

	return image, nil
}

func (d *DB) Delete(ctx context.Context, id string) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM images WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	return nil
}

func (d *DB) Update(ctx context.Context, id string) error {
	query := `UPDATE images SET status = $1 WHERE id = $2`
	_, err := d.db.ExecContext(ctx, query, "completed", id)
	return err
}

func (d *DB) UpdatePathsAndStatus(ctx context.Context, id, thumb, watermark, resized string) error {
	query := `
		UPDATE images
		SET thumbnail_path = $1,
		    watermark_path = $2,
		    resized_path = $3,
		    status = $4
		WHERE id = $5
	`

	_, err := d.db.ExecContext(ctx, query, thumb, watermark, resized, "completed", id)
	if err != nil {
		return fmt.Errorf("update paths and status: %w", err)
	}

	return nil
}
