CREATE TABLE IF NOT EXISTS images (
    id TEXT PRIMARY KEY,
    original_path TEXT NOT NULL,
    thumbnail_path TEXT,
    watermark_path TEXT,
    resized_path TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);