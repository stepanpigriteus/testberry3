CREATE TABLE IF NOT EXISTS images (
    id TEXT PRIMARY KEY,
    original_path TEXT NOT NULL,
    processed_path TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);