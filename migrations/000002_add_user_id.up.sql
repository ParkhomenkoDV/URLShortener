-- +migrate Up
ALTER TABLE urls ADD COLUMN user_id VARCHAR(36);
