-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL SET DEFAULT "unset";

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;