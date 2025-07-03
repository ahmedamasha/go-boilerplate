-- +goose Up
ALTER TABLE users ADD COLUMN password VARCHAR(255) NOT NULL DEFAULT '';


-- +goose Down
ALTER TABLE users DROP COLUMN password;