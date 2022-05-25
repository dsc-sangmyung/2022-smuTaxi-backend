-- name: CreateUser :one
INSERT INTO users (
    id,
    name,
    email,
    gender
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GerUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: EnterRoom :one
UPDATE users
SET room_id = $2
WHERE id = $1
RETURNING *;

-- name: QuitRoom :one
UPDATE users
SET room_id = NULL
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
