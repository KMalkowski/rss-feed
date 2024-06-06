-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.description, posts.published_at, posts.feed_id from posts
LEFT JOIN feeds ON feeds.id = posts.feed_id
WHERE feeds.user_id = $1
LIMIT $2;
