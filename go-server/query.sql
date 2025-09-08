-- query file for sqlc

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

-- name: GetLatestCapacity :one
SELECT *
FROM gym_capacities
WHERE location_id = $1
ORDER BY last_updated_at DESC
LIMIT 1;

-- name: CreateGymCapacity :one
INSERT INTO gym_capacities (
  location_id, percentage, last_updated_at
) VALUES (
  $1, $2, $3
) 
RETURNING *;

-- name: GetGymCapacitiesBetween :many
SELECT *
FROM gym_capacities
WHERE last_updated_at >= $1
  AND last_updated_at <= $2
ORDER BY last_updated_at;
