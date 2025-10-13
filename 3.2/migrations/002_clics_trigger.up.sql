CREATE OR REPLACE FUNCTION increment_click_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE short_links
    SET click_count = click_count + 1
    WHERE id = NEW.short_link_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_increment_click_count
AFTER INSERT ON link_visits
FOR EACH ROW
EXECUTE FUNCTION increment_click_count();