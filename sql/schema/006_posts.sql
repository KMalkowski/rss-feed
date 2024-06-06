-- +goose Up
CREATE TABLE posts (id UUID PRIMARY KEY, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP NOT NULL,
title TEXT NOT NULL, url TEXT NOT NULL, description TEXT, published_at TIMESTAMP NOT NULL, feed_id UUID NOT NULL,
CONSTRAINT unique_url UNIQUE(url));

-- +goose Down
DROP TABLE posts;
