-- name: GetDiningUser :one
SELECT * FROM dining_users
WHERE device_id = $1 LIMIT 1;

-- name: GetDiningUserAll :many
SELECT * FROM dining_users
WHERE user_id = $1;

-- name: CreateDiningUser :one
INSERT INTO dining_users (
  user_id, device_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateDiningUserSession :exec
UPDATE dining_users
  SET last_session_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteDiningUser :exec
DELETE FROM dining_users
WHERE user_id = $1;
