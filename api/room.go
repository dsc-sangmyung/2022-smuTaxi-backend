package api

import (
	"database/sql"
	db "smutaxi/db/sqlc"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Create Room
type createRoomReq struct {
	Source      string `json:"source" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Date        string `json:"date" validate:"required"`
	Time        string `json:"time" validate:"required"`
}

func (server *Server) createRoom(ctx *fiber.Ctx) error {
	req := new(createRoomReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	stringToDate, _ := time.Parse("2006-01-02", req.Date)
	stringToTime, _ := time.Parse("15:04", req.Time)

	arg := db.CreateRoomParams{
		Source:      req.Source,
		Destination: req.Destination,
		Date:        stringToDate,
		Time:        stringToTime,
	}

	room, err := server.store.CreateRoom(ctx.Context(), arg)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(room)
}

type People struct {
	Name   string   `json:"name"`
	Gender string   `json:"gender"`
	Items  []string `json:"item"`
}

type getRoomRes struct {
	Room db.Room `json:"room"`
	// Items   [][]string       `json:"item"`
	Peoples []db.FindUserRow `json:"people"`
}

// Get Room by room_id
func (server *Server) getRoom(ctx *fiber.Ctx) error {

	result := new(getRoomRes)

	roomID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	room, err := server.store.GetRoom(ctx.Context(), roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	result.Room = room

	member := room.Member

	for _, v := range member {
		people, err := server.store.FindUser(ctx.Context(), v)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		}

		result.Peoples = append(result.Peoples, people)
	}

	return ctx.JSON(result)
}

func (server *Server) getAllRooms(ctx *fiber.Ctx) error {
	rooms, err := server.store.GetRooms(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(rooms)
}

type findRoomsReq struct {
	Source      string `json:"source" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Date        string `json:"date" validate:"required"`
	Time        string `json:"time" validate:"required"`
}

func (server *Server) findRooms(ctx *fiber.Ctx) error {
	req := new(findRoomsReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	stringToDate, _ := time.Parse("2006-01-02", req.Date)
	stringToTime, _ := time.Parse("15:04", req.Time)

	arg := db.FindRoomsParams{
		Source:      req.Source,
		Destination: req.Destination,
		Date:        stringToDate,
		Time:        stringToTime,
	}

	rooms, err := server.store.FindRooms(ctx.Context(), arg)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(rooms)
}

// Add member to room
type addMemberReq struct {
	Member string `json:"member" validate:"required"`
}

func (server *Server) addMember(ctx *fiber.Ctx) error {
	req := new(addMemberReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	roomID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.AddMemberParams{
		RoomID:      roomID,
		ArrayAppend: req.Member,
	}

	room, err := server.store.AddMember(ctx.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(room)
}

func (server *Server) deleteRoom(ctx *fiber.Ctx) error {
	roomID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	err = server.store.DeleteRoom(ctx.Context(), roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(nil)
}
