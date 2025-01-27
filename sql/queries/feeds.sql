-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextToFetch :many
SELECT * FROM feeds ORDER BY last_fetched_at LIMIT $1;

-- name: MarkFetched :one
UPDATE feeds SET last_fetched_at = $1, updated_at = $1
WHERE id = $2
RETURNING *;
