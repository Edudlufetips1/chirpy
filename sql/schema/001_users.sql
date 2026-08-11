-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;

-- +goose Up
ALTER TABLE users 
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users 
DROP COLUMN is_chirpy_red;