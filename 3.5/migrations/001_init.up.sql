CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    total_seats INTEGER NOT NULL CHECK (total_seats > 0),
    available_seats INTEGER NOT NULL CHECK (available_seats >= 0),  
    requires_payment BOOLEAN DEFAULT true,
    booking_ttl INTEGER NOT NULL DEFAULT 3600 CHECK (booking_ttl > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_date ON events (date);

CREATE INDEX idx_events_available_seats ON events (available_seats);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    telegram_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_role ON users (role);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    event_id UUID NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'confirmed', 'cancelled', 'expired')
    ),
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP,
    CONSTRAINT chk_confirmed_at CHECK (
        (
            status = 'confirmed' AND confirmed_at IS NOT NULL
        )
        OR (
            status != 'confirmed' AND confirmed_at IS NULL
        )
    )
);


CREATE INDEX idx_bookings_status_expires ON bookings (status, expires_at);

CREATE INDEX idx_bookings_event_id ON bookings (event_id);

CREATE INDEX idx_bookings_user_id ON bookings (user_id);

CREATE INDEX idx_bookings_expires_at ON bookings (expires_at)
WHERE
    status = 'pending';

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    booking_id UUID NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    user_email VARCHAR(255) NOT NULL,
    telegram_id VARCHAR(100),
    type VARCHAR(50) NOT NULL CHECK (
        type IN (
            'booking_created',
            'booking_expiring',
            'booking_cancelled',
            'booking_confirmed'
        )
    ),
    message TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'sent', 'failed')
    ),
    created_at TIMESTAMP DEFAULT NOW(),
    sent_at TIMESTAMP
);

CREATE INDEX idx_notifications_status ON notifications (status);

CREATE INDEX idx_notifications_booking_id ON notifications (booking_id);

CREATE INDEX idx_notifications_created_at ON notifications (created_at)
WHERE
    status = 'pending';

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();