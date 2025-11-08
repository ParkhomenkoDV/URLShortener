-- +migrate Up
ALTER TABLE urls ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

-- Добавляем индекс для ускорения поиска неудаленных URL
CREATE INDEX idx_urls_is_deleted ON urls(is_deleted);
