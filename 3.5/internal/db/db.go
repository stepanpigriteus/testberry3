package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"threeFive/domain"

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

func (d *DB) Create(ctx context.Context, event domain.Event) (string, error) {
	id := uuid.New()
	fmt.Println(event.TotalSeats)
	query := `
		INSERT INTO events (
			id, title, description, date,
			total_seats, available_seats,
			requires_payment, booking_ttl
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`
	ttl := time.Duration(event.BookingTTL)
	if ttl <= 0 {
		ttl = 3600
	}

	var newID string
	err := d.db.Master.QueryRowContext(ctx, query,
		id,
		event.Title,
		event.Description,
		event.Date,
		event.TotalSeats,
		event.AvailableSeats,
		event.RequiresPayment,
		int(ttl),
	).Scan(&newID)
	if err != nil {
		return "", err
	}

	return newID, nil

}

func (d *DB) Book(ctx context.Context, id string) (string, error) {
	tx, err := d.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	var availableSeats int
	var bookingTTL int

	err = tx.QueryRowContext(ctx,
		`SELECT available_seats, booking_ttl FROM events WHERE id = $1 FOR UPDATE`,
		id,
	).Scan(&availableSeats, &bookingTTL)
	if err == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("select event: %w", err)
	}
	if availableSeats <= 0 {
		return "", domain.ErrInvalidSeats
	}

	var userID string
	err = tx.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, "test_email").Scan(&userID)
	if err == sql.ErrNoRows {
		userID = uuid.New().String()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO users (id, email, name, role)
			VALUES ($1, $2, $3, 'user')
		`, userID, "test_email", "test_name")
		if err != nil {
			return "", fmt.Errorf("insert user: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("select user: %w", err)
	}

	bookingID := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(bookingTTL) * time.Second)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO bookings (id, event_id, user_email, user_name, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		bookingID,
		id,
		"user@example.com",
		"Anonymous",
		"pending",
		expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("insert booking: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE events
		SET available_seats = available_seats - 1
		WHERE id = $1
	`, id)
	if err != nil {
		return "", fmt.Errorf("update seats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("commit tx: %w", err)
	}

	return bookingID, nil
}

func (d *DB) GetEvent(ctx context.Context, id string) (domain.Event, error) {
	var event domain.Event
	return event, nil
}

func (d *DB) Update(ctx context.Context, id string) error {
	return nil
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