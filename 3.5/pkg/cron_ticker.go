package pkg

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func StartBookingCleaner(ctx context.Context, db *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Booking cleaner started")
	cleanExpiredBookings(ctx, db)

	for {
		select {
		case <-ticker.C:
			cleanExpiredBookings(ctx, db)
		case <-ctx.Done():
			log.Println("Booking cleaner stopped")
			return
		}
	}
}

func cleanExpiredBookings(ctx context.Context, db *sql.DB) {
	query := `
		WITH expired AS (
			UPDATE bookings
			SET status = 'expired'
			WHERE status = 'pending'
			AND expires_at <= NOW()
			RETURNING event_id
		)
		UPDATE events
		SET available_seats = LEAST(available_seats + 1, total_seats)
		FROM expired
		WHERE events.id = expired.event_id
	`

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error cleaning expired bookings: %v", err)
		return
	}

	if rows, _ := result.RowsAffected(); rows > 0 {
		log.Printf("Cleaned %d expired bookings", rows)
	}
}
