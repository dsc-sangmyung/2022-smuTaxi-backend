package api

import (
	"database/sql"
	db "smutaxi/db/sqlc"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EnterRoomReq struct {
	UserID string
	RoomID int64
}

func (server *Server) EnterRoom(ctx *fiber.Ctx) error {
	req := new(EnterRoomReq)

	if err := ctx.BodyParser(req); err != nil {
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
		},
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

func (server *Server) QuitRoom(ctx *fiber.Ctx) error {

	return nil
}
