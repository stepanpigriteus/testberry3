DROP TRIGGER IF EXISTS trg_comments_tsv ON comments;
DROP FUNCTION IF EXISTS update_tsv();

DROP INDEX IF EXISTS idx_comments_tsv;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_created_at;

DROP TABLE IF EXISTS comments;