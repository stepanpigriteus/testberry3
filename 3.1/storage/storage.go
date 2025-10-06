package storage

import (
	"context"
	"database/sql"
	"fmt"
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type StorageImpl struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewStorage(ctx context.Context, masterDSN string, slaveDSNs []string, logger zerolog.Logger) *StorageImpl {
	fmt.Println(masterDSN)
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("init database error %s", err)
	}
	return &StorageImpl{
		db:     db,
		logger: logger,
	}
}

func (st *StorageImpl) CreateNotify(ctx context.Context, notify domain.Notify) error {
	st.logger.Info().Any("Создаётся уведомление: ", notify)

	if err := st.db.Master.PingContext(ctx); err != nil {
		st.logger.Err(err).Msgf("Ошибка соединения с БД: %v\n", err)
		return fmt.Errorf("нет соединения с БД: %w", err)
	}
	query := `
		INSERT INTO notify (id, timing, descript, status, retry, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := st.db.ExecContext(ctx, query,
		notify.Id,
		notify.Timing,
		notify.Descript,
		notify.Status,
		notify.Retry,
		notify.CreatedAt,
	)

	if err != nil {
		st.logger.Err(err).Msgf("Ошибка запроса: %v\nЗапрос: %s\nПараметры: %+v\n", err, query, notify)
		return err
	}

	return nil
}

func (st *StorageImpl) GetNotify(ctx context.Context, id string) (domain.Notify, error) {
	var n domain.Notify
	if err := st.db.Master.PingContext(ctx); err != nil {
		st.logger.Err(err).Msgf("Ошибка соединения с БД: %v\n", err)
		return n, fmt.Errorf("нет соединения с БД: %w", err)
	}
	query := `
		SELECT id, timing, descript, status, retry, created_at
		FROM notify
		WHERE id = $1
	`

	err := st.db.QueryRowContext(ctx, query, id).Scan(
		&n.Id,
		&n.Timing,
		&n.Descript,
		&n.Status,
		&n.Retry,
		&n.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return n, err
		}
		return n, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	return n, nil
}

func (st *StorageImpl) DeleteNotify(ctx context.Context, id string) error {
	if err := st.db.Master.PingContext(ctx); err != nil {
		st.logger.Err(err).Msgf("Ошибка соединения с БД: %v\n", err)
		return fmt.Errorf("нет соединения с БД: %w", err)
	}

	query := `
		UPDATE notify
		SET status = 'cancelled'
		WHERE id = $1
	`

	result, err := st.db.ExecContext(ctx, query, id)
	if err != nil {
		st.logger.Err(err).Msgf("Ошибка обновления статуса уведомления id=%s", id)
		return fmt.Errorf("ошибка обновления: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("уведомление с id=%s не найдено", id)
	}

	return nil
}
