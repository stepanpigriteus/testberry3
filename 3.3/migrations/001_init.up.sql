CREATE TABLE comments (
    id SERIAL PRIMARY KEY,                  -- Уникальный ID комментария (автоинкремент)
    parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,  -- ID родительского комментария (NULL для корневых)
    text TEXT NOT NULL,                     -- Текст комментария
    author VARCHAR(255) NOT NULL,           -- Автор (username или email)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Дата создания
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Дата обновления
    deleted_at TIMESTAMP WITH TIME ZONE     -- Мягкое удаление
);

-- Индексы для производительности
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- Для полнотекстового поиска (в PostgreSQL)
ALTER TABLE comments ADD COLUMN tsv TSVECTOR;
CREATE INDEX idx_comments_tsv ON comments USING GIN(tsv);

-- Триггер для обновления tsv при вставке/обновлении
CREATE OR REPLACE FUNCTION update_tsv() RETURNS TRIGGER AS $$
BEGIN
    -- Многоязычный поиск: русский + английский
    NEW.tsv :=
        to_tsvector('russian', NEW.text) ||
        to_tsvector('english', NEW.text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_comments_tsv
BEFORE INSERT OR UPDATE ON comments
FOR EACH ROW
EXECUTE PROCEDURE update_tsv();
