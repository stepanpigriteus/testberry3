CREATE TABLE IF NOT EXISTS short_links (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(32) UNIQUE NOT NULL,    
    original_url TEXT NOT NULL,                             
    created_at TIMESTAMP DEFAULT NOW(),                 
    click_count INT DEFAULT 0         
);

CREATE INDEX IF NOT EXISTS idx_short_links_code
    ON short_links (short_code);


CREATE TABLE IF NOT EXISTS link_visits (
    id SERIAL PRIMARY KEY,
    short_link_id INT NOT NULL REFERENCES short_links(id) ON DELETE CASCADE,
    visited_at TIMESTAMP DEFAULT NOW(),
    user_agent TEXT,
    ip_address INET,
    device_type VARCHAR(32)
);

CREATE INDEX IF NOT EXISTS idx_link_visits_link_id
    ON link_visits (short_link_id);

CREATE INDEX IF NOT EXISTS idx_link_visits_visited_at
    ON link_visits (visited_at);