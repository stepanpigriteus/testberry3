package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"threeSixth/domain"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (d *DB) Create(ctx context.Context, item *domain.Item) (domain.Item, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO items (id, type, category, amount, date, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := d.db.QueryRowContext(ctx, query,
		id,
		item.Type,
		item.Category,
		item.Amount,
		item.Date,
		item.Description,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return *item, domain.ErrDuplicateKey
		}
		return *item, fmt.Errorf("insert item: %w", err)
	}

	return *item, nil
}

func (d *DB) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Item, error) {
	query := `
		SELECT id, type, category, amount, date, description, created_at, updated_at
		FROM items
		WHERE 1=1
	`
	args := []any{}
	argIdx := 1

	if filter.From != "" {
		query += fmt.Sprintf(" AND date >= $%d", argIdx)
		args = append(args, filter.From)
		argIdx++
	}
	if filter.To != "" {
		query += fmt.Sprintf(" AND date <= $%d", argIdx)
		args = append(args, filter.To)
		argIdx++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, filter.Category)

	}

	allowedSortBy := map[string]bool{
		"date":       true,
		"amount":     true,
		"category":   true,
		"created_at": true,
	}

	order := "ASC"
	if strings.ToUpper(filter.Order) == "DESC" {
		order = "DESC"
	}

	if filter.SortBy != "" {
		if allowedSortBy[filter.SortBy] {
			query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, order)
		} else {
			return nil, fmt.Errorf("invalid sort field: %s", filter.SortBy)
		}
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Category,
			&item.Amount,
			&item.Date,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return items, nil
}

func (d *DB) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	return nil, nil
}

func (d *DB) Update(ctx context.Context, item *domain.Item) (domain.Item, error) {
	query := `
		UPDATE items
		SET 
			type = $1,
			category = $2,
			amount = $3,
			date = $4,
			description = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, type, category, amount, date, description, created_at, updated_at
	`

	var updated domain.Item
	err := d.db.QueryRowContext(ctx, query,
		item.Type,
		item.Category,
		item.Amount,
		item.Date,
		item.Description,
		item.ID,
	).Scan(
		&updated.ID,
		&updated.Type,
		&updated.Category,
		&updated.Amount,
		&updated.Date,
		&updated.Description,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Item{}, fmt.Errorf("no item found with id %s", item.ID)
		}
		return domain.Item{}, fmt.Errorf("update query: %w", err)
	}

	return updated, nil
}

func (d *DB) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (d *DB) GetAnalytics(ctx context.Context, filter domain.AnalyticsFilter) (*domain.AnalyticsResult, error) {
	if filter.From != "" {
		if _, err := time.Parse("2006-01-02", filter.From); err != nil {
			return nil, fmt.Errorf("invalid From date: %w", err)
		}
	}
	if filter.To != "" {
		if _, err := time.Parse("2006-01-02", filter.To); err != nil {
			return nil, fmt.Errorf("invalid To date: %w", err)
		}
	}

	order := "ASC"
	allowedGroupBy := map[string]bool{
		"type":     true,
		"category": true,
		"date":     true,
	}
	if filter.GroupBy != "" && !allowedGroupBy[filter.GroupBy] {
		return nil, fmt.Errorf("invalid group_by field: %s", filter.GroupBy)
	}

	args := []any{}
	queryTotal := `
        SELECT
            COUNT(*) AS count,
            COALESCE(SUM(amount), 0) AS sum,
            COALESCE(AVG(amount), 0) AS avg,
            PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) AS median,
            PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) AS percentile_90
        FROM items
        WHERE 1=1
    `

	if filter.From != "" {
		queryTotal += fmt.Sprintf(" AND date >= $%d", len(args)+1)
		args = append(args, filter.From)
	}
	if filter.To != "" {
		queryTotal += fmt.Sprintf(" AND date <= $%d", len(args)+1)
		args = append(args, filter.To)
	}
	if filter.Type != "" {
		queryTotal += fmt.Sprintf(" AND type = $%d", len(args)+1)
		args = append(args, filter.Type)
	}

	var result domain.AnalyticsResult
	if err := d.db.QueryRowContext(ctx, queryTotal, args...).Scan(
		&result.Count,
		&result.Sum,
		&result.Avg,
		&result.Median,
		&result.Percentile90,
	); err != nil {
		return nil, fmt.Errorf("total analytics query: %w", err)
	}

	if filter.GroupBy != "" {
		groupArgs := []any{}
		queryGrouped := fmt.Sprintf(`
            SELECT %s AS group_name,
                COUNT(*) AS count,
                COALESCE(SUM(amount), 0) AS sum,
                COALESCE(AVG(amount), 0) AS avg,
                PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) AS median,
                PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) AS percentile_90
            FROM items
            WHERE 1=1
        `, pq.QuoteIdentifier(filter.GroupBy))

		if filter.From != "" {
			queryGrouped += fmt.Sprintf(" AND date >= $%d", len(groupArgs)+1)
			groupArgs = append(groupArgs, filter.From)
		}
		if filter.To != "" {
			queryGrouped += fmt.Sprintf(" AND date <= $%d", len(groupArgs)+1)
			groupArgs = append(groupArgs, filter.To)
		}
		if filter.Type != "" {
			queryGrouped += fmt.Sprintf(" AND type = $%d", len(groupArgs)+1)
			groupArgs = append(groupArgs, filter.Type)
		}

		queryGrouped += fmt.Sprintf(" GROUP BY %s ORDER BY %s %s", filter.GroupBy, filter.GroupBy, order)

		rows, err := d.db.QueryContext(ctx, queryGrouped, groupArgs...)
		if err != nil {
			return nil, fmt.Errorf("grouped analytics query: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var g domain.GroupedMetrics
			if err := rows.Scan(&g.Group, &g.Count, &g.Sum, &g.Avg, &g.Median, &g.Percentile90); err != nil {
				return nil, fmt.Errorf("scan grouped metrics: %w", err)
			}
			result.Grouped = append(result.Grouped, g)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}
	}

	if result.Sum < 0 {
		return nil, fmt.Errorf("result sum < 0 ")
	}

	return &result, nil
}

func (d *DB) Close() error {
	if d.db == nil {
		return nil
	}

	d.logger.Info().Msg("Closing database connection...")
	if err := d.db.Master.Close(); err != nil {
		d.logger.Error().Err(err).Msg("Failed to close database connection")
		return err
	}

	d.logger.Info().Msg("Database connection closed successfully")
	return nil
}

func (d *DB) GetMaster() *sql.DB {
	return d.db.Master
}
