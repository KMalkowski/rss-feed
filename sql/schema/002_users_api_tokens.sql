-- +goose Up
ALTER TABLE users
ADD COLUMN api_key VARCHAR(64) NOT NULL DEFAULT encode(sha256(random()::text::bytea), 'hex'),
ADD CONSTRAINT unique_api_key UNIQUE (api_key);

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key,
DROP CONSTRAINT unique_api_key;
