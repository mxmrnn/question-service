-- +goose Up
CREATE TABLE IF NOT EXISTS questions (
    id         SERIAL PRIMARY KEY,
    text       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS questions;