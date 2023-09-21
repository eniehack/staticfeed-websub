-- name: InsertNewToken :exec
INSERT INTO users (id, access_token, refresh_token, created_at, expired_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateToken :exec
UPDATE users SET access_token = ?, refresh_token = ?, expired_at = ?, updated_at = ?
WHERE id = ?;

-- name: GetUserFromID :one
SELECT * FROM users WHERE id = ?;
