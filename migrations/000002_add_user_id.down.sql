-- +migrate Down
ALTER TABLE urls DROP COLUMN user_id;
