CREATE TABLE IF NOT EXISTS notify (
    id SERIAL PRIMARY KEY,
    timing TIMESTAMPTZ NOT NULL,
    descript TEXT,
    status TEXT,
    retry INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);