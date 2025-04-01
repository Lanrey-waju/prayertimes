-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER cleanup_old_prayer_times 
AFTER INSERT ON prayer_times
BEGIN
    DELETE FROM prayer_times 
    WHERE DATE(created_at) < DATE('now', '-10');
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER cleanup_old_prayer_times;
-- +goose StatementEnd
