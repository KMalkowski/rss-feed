-- name: CreateFollow :one
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUsersFollows :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: DeleteFollow :execrows
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2;
