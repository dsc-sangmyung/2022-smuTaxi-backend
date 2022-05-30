-- name: CreateRoom :one
INSERT INTO room (
    source,
    destination,
    date,
    time
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetRoom :one
SELECT * FROM room
WHERE room_id = $1 LIMIT 1;

-- name: GetRooms :many
SELECT * FROM room;

-- name: ListTodayRoom :many
SELECT * FROM room
WHERE date = CURRENT_DATE;

-- name: AddMember :one
UPDATE room
SET member = array_append(member, $2)
WHERE room_id = $1
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM room
WHERE room_id = $1;
