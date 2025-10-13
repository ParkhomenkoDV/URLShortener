-- +migrate Up
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_urls_original_url ON urls(original_url);
CREATE UNIQUE INDEX idx_urls_short_url ON urls(short_url);