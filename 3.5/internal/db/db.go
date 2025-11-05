package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func (d *DB) Book(ctx context.Context, eventID string, userID string) (string, error) {
	tx, err := d.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var dummy string
	err = tx.QueryRowContext(ctx, `SELECT id FROM users WHERE id = $1`, userID).Scan(&dummy)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user does not exist")
	}
	if err != nil {
		return "", fmt.Errorf("select user: %w", err)
	}

	var availableSeats int
	var bookingTTL int
	err = tx.QueryRowContext(ctx, `
		SELECT available_seats, booking_ttl
		FROM events
		WHERE id = $1
		FOR UPDATE
	`, eventID).Scan(&availableSeats, &bookingTTL)
	if err == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("select event: %w", err)
	}

	if availableSeats <= 0 {
		return "", domain.ErrInvalidSeats
	}

	bookingID := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(bookingTTL) * time.Second)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO bookings (id, event_id, user_id, status, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`,
		bookingID,
		eventID,
		userID,
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
	`, eventID)
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
	query := `SELECT id, title, description, date, total_seats, available_seats, requires_payment, booking_ttl, created_at
	          FROM events WHERE id = $1`

	err := d.db.Master.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.RequiresPayment,
		&event.BookingTTL,
		&event.CreatedAt,
	)
	if err != nil {
		return domain.Event{}, err
	}
	fmt.Println(event)
	return event, nil
}

func (d *DB) Update(ctx context.Context, eventId, bookId string) error {
	tx, err := d.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRowContext(ctx, `
		SELECT status FROM bookings
		WHERE event_id = $1 AND id = $2
		FOR UPDATE
	`, eventId, bookId).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("booking not found")
		}
		return fmt.Errorf("check booking: %w", err)
	}

	if status != "pending" {
		return fmt.Errorf("cannot confirm booking with status=%q", status)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE bookings
		SET status = 'confirmed', confirmed_at = NOW()
		WHERE event_id = $1 AND id = $2
	`, eventId, bookId)
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (d *DB) CreateUser(ctx context.Context, user domain.User) (string, error) {
	userId := uuid.New().String()
	query := `
		INSERT INTO users (id, email, name, role, telegram_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := d.db.ExecContext(ctx, query,
		userId,
		user.Email,
		user.Name,
		user.Role,
		user.TelegramID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", fmt.Errorf("user with email %s already exists", user.Email)
		}
		return "", fmt.Errorf("insert user: %w", err)
	}

	return string(userId), nil
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
