package domain

import "time"

type Event struct {
	ID              string        `db:"id"`
	Title           string        `db:"title"`
	Description     string        `db:"description"`
	Date            time.Time     `db:"date"`
	TotalSeats      int           `db:"total_seats"`
	AvailableSeats  int           `db:"available_seats"`
	RequiresPayment bool          `db:"requires_payment"`
	BookingTTL      time.Duration `db:"booking_ttl"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
}

type Booking struct {
	ID          string        `db:"id"`
	EventID     string        `db:"event_id"`
	UserID      *string       `db:"user_id"`
	UserEmail   string        `db:"user_email"`
	UserName    string        `db:"user_name"`
	Status      BookingStatus `db:"status"`
	CreatedAt   time.Time     `db:"created_at"`
	ExpiresAt   time.Time     `db:"expires_at"`
	ConfirmedAt *time.Time    `db:"confirmed_at"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusExpired   BookingStatus = "expired"
)

type User struct {
	ID         string    `db:"id"`
	Email      string    `db:"email"`
	Name       string    `db:"name"`
	Role       UserRole  `db:"role"`
	TelegramID string    `db:"telegram_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type Notification struct {
	ID         string             `db:"id"`
	BookingID  string             `db:"booking_id"`
	UserEmail  string             `db:"user_email"`
	TelegramID string             `db:"telegram_id"`
	Type       NotificationType   `db:"type"`
	Message    string             `db:"message"`
	Status     NotificationStatus `db:"status"`
	CreatedAt  time.Time          `db:"created_at"`
	SentAt     *time.Time         `db:"sent_at"`
}

type NotificationType string

const (
	NotificationBookingCreated   NotificationType = "booking_created"
	NotificationBookingExpiring  NotificationType = "booking_expiring"
	NotificationBookingCancelled NotificationType = "booking_cancelled"
	NotificationBookingConfirmed NotificationType = "booking_confirmed"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)
