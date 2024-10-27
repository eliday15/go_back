-- name: CreateFollow :one
INSERT INTO follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFollowsByUserID :many
SELECT * FROM follows WHERE user_id = $1;

-- name: DeleteFollow :exec
DELETE FROM follows WHERE id = $1 AND user_id = $2;