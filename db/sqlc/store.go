package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type EnterRoomTxParams struct {
	UserID string
	RoomID sql.NullInt64
	Item   []string
}

type EnterRoomTxResult struct {
	User User
	Room Room
}

type QuitRoomTxParams struct {
	UserID string
	RoomID sql.NullInt64
}

type QuitRoomTxResult struct {
	User User
	Room Room
}

func (store *Store) QuitRoomTx(ctx context.Context, arg QuitRoomTxParams) (QuitRoomTxResult, error) {
	var result QuitRoomTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 0. User의 Roomid 초기화
		err = q.QuitRoom(ctx, arg.UserID)
		if err != nil {
			return err
		}

		// 1. Room의 member에서 User 삭제
		err = q.RemoveMember(ctx, RemoveMemberParams{
			RoomID:      arg.RoomID.Int64,
			ArrayRemove: arg.UserID,
		})

		if err != nil {
			return err
		}

		return err
	})
	return result, err
}

func (store *Store) EnterRoomTx(ctx context.Context, arg EnterRoomTxParams) (EnterRoomTxResult, error) {
	var result EnterRoomTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 0. User가 현재 들어가있는 방이 없는지 체크한다.
		result.User, err = q.GetUser(ctx, arg.UserID)
		if result.User.RoomID.Valid {
			return fmt.Errorf("user is already in room")
		}
		if err != nil {
			return err
		}

		// 1. Room이 존재하는지 체크한다.
		result.Room, err = q.GetRoom(ctx, arg.RoomID.Int64)
		if err != nil {
			return err
		}

		// 1. room_id의 인원이 4명인지 체크한다. 4명이면 거부.
		if len(result.Room.Member) == 4 {
			return fmt.Errorf("room is full")
		}
		// 2. user의 item 정보를 갱신해줌.
		store.SetItems(ctx, SetItemsParams{
			ID:   arg.UserID,
			Item: arg.Item,
		})

		// 3. 4명이 아니라면 유저의 room_id값을 갱신해줌.
		result.User, err = store.EnterRoom(ctx, EnterRoomParams{
			ID:     arg.UserID,
			RoomID: arg.RoomID,
		})
		if err != nil {
			return err
		}

		// 4. room_id에 해당하는 room의 member에 user를 추가.
		result.Room, err = store.AddMember(ctx, AddMemberParams{
			RoomID:      arg.RoomID.Int64,
			ArrayAppend: arg.UserID,
		})
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}
