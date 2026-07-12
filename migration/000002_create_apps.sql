-- +goose Up
CREATE TABLE IF NOT EXISTS apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_apps_name ON apps(name);

-- +goose Down
DROP INDEX IF EXISTS idx_apps_name;
DROP TABLE IF EXISTS apps;
