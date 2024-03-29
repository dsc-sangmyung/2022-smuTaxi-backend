// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    name,
    email,
    gender
) VALUES (
    $1, $2, $3, $4
) RETURNING id, name, email, gender, item, created_at, room_id
`

type CreateUserParams struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Gender,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Gender,
		pq.Array(&i.Item),
		&i.CreatedAt,
		&i.RoomID,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const enterRoom = `-- name: EnterRoom :one
UPDATE users
SET room_id = $2
WHERE id = $1
RETURNING id, name, email, gender, item, created_at, room_id
`

type EnterRoomParams struct {
	ID     string        `json:"id"`
	RoomID sql.NullInt64 `json:"room_id"`
}

func (q *Queries) EnterRoom(ctx context.Context, arg EnterRoomParams) (User, error) {
	row := q.db.QueryRowContext(ctx, enterRoom, arg.ID, arg.RoomID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Gender,
		pq.Array(&i.Item),
		&i.CreatedAt,
		&i.RoomID,
	)
	return i, err
}

const findUser = `-- name: FindUser :one
SELECT name, gender, item FROM users
WHERE id = $1 LIMIT 1
`

type FindUserRow struct {
	Name   string   `json:"name"`
	Gender string   `json:"gender"`
	Item   []string `json:"item"`
}

func (q *Queries) FindUser(ctx context.Context, id string) (FindUserRow, error) {
	row := q.db.QueryRowContext(ctx, findUser, id)
	var i FindUserRow
	err := row.Scan(&i.Name, &i.Gender, pq.Array(&i.Item))
	return i, err
}

const gerUserForUpdate = `-- name: GerUserForUpdate :one
SELECT id, name, email, gender, item, created_at, room_id FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GerUserForUpdate(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, gerUserForUpdate, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Gender,
		pq.Array(&i.Item),
		&i.CreatedAt,
		&i.RoomID,
	)
	return i, err
}

const getItems = `-- name: GetItems :one
SELECT item from users
WHERE id = $1
`

func (q *Queries) GetItems(ctx context.Context, id string) ([]string, error) {
	row := q.db.QueryRowContext(ctx, getItems, id)
	var item []string
	err := row.Scan(pq.Array(&item))
	return item, err
}

const getUser = `-- name: GetUser :one
SELECT id, name, email, gender, item, created_at, room_id FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Gender,
		pq.Array(&i.Item),
		&i.CreatedAt,
		&i.RoomID,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, name, email, gender, item, created_at, room_id FROM users
ORDER BY created_at
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Gender,
			pq.Array(&i.Item),
			&i.CreatedAt,
			&i.RoomID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const quitRoom = `-- name: QuitRoom :exec
UPDATE users
SET room_id = NULL
WHERE id = $1
`

func (q *Queries) QuitRoom(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, quitRoom, id)
	return err
}

const setItems = `-- name: SetItems :exec
UPDATE users
SET item = $2
WHERE id = $1
`

type SetItemsParams struct {
	ID   string   `json:"id"`
	Item []string `json:"item"`
}

func (q *Queries) SetItems(ctx context.Context, arg SetItemsParams) error {
	_, err := q.db.ExecContext(ctx, setItems, arg.ID, pq.Array(arg.Item))
	return err
}
