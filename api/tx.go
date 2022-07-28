package api

import (
	"database/sql"
	"fmt"
	db "smutaxi/db/sqlc"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EnterRoomReq struct {
	UserID string   `json:"user_id" validate:"required"`
	RoomID int64    `json:"room_id" validate:"required"`
	Item   []string `json:"item"`
}

func (server *Server) EnterRoom(ctx *fiber.Ctx) error {
	req := new(EnterRoomReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	fmt.Println(req)

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	roomID, err := strconv.ParseInt(ctx.Params("roomID"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.EnterRoomTxParams{
		UserID: req.UserID,
		RoomID: sql.NullInt64{
			Int64: roomID,
			Valid: true,
		},
		Item: req.Item,
	}

	res, err := server.store.EnterRoomTx(ctx.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(res)
}

type QuitRoomReq struct {
	UserID string `json:"user_id" validate:"required"`
	RoomID int64  `json:"room_id" validate:"required"`
}

func (server *Server) QuitRoom(ctx *fiber.Ctx) error {
	req := new(QuitRoomReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	fmt.Println("REQ: ", req)

	arg := db.QuitRoomTxParams{
		UserID: req.UserID,
		RoomID: sql.NullInt64{
			Int64: req.RoomID,
			Valid: true,
		},
	}

	res, err := server.store.QuitRoomTx(ctx.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(res)
}
