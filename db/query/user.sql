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

-- name: FindUser :one
SELECT name, gender, item FROM users
WHERE id = $1 LIMIT 1;

-- name: GerUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: SetItems :exec
UPDATE users
SET item = $2
WHERE id = $1;

-- name: GetItems :one
SELECT item from users
WHERE id = $1;

-- name: EnterRoom :one
UPDATE users
SET room_id = $2
WHERE id = $1
RETURNING *;

-- name: QuitRoom :exec
UPDATE users
SET room_id = NULL
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
