package api

import (
	"database/sql"
	db "smutaxi/db/sqlc"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Create User
type createUserReq struct {
	ID     string `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
}

func (server *Server) createUser(ctx *fiber.Ctx) error {
	req := new(createUserReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.CreateUserParams{
		ID:     req.ID,
		Name:   req.Name,
		Email:  req.Email,
		Gender: req.Gender,
	}

	user, err := server.store.CreateUser(ctx.Context(), arg)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(user)
}

// Get User by id
type getUserReq struct {
	ID string
}

func (server *Server) getUser(ctx *fiber.Ctx) error {

	req := new(getUserReq)

	req.ID = ctx.Params("id")
	if req.ID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty id field",
		})
	}

	user, err := server.store.GetUser(ctx.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(user)
}

// Get All Users
func (server *Server) getUsers(ctx *fiber.Ctx) error {

	users, err := server.store.GetUsers(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}
	return ctx.JSON(users)
}

// Modify User
func (server *Server) patchUser() {}

type setItemReq struct {
	Item []string `json:"item" validate:"required"`
}

func (server *Server) setItem(ctx *fiber.Ctx) error {
	req := new(setItemReq)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.SetItemsParams{
		ID:   ctx.Params("id"),
		Item: req.Item,
	}

	err := server.store.SetItems(ctx.Context(), arg)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(fiber.Map{
		"message": "set item successful.",
	})

}

func (server *Server) deleteUser(ctx *fiber.Ctx) error {
	idReq := ctx.Params("id")
	if idReq == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty id field",
		})
	}

	err := server.store.DeleteUser(ctx.Context(), idReq)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return ctx.JSON(fiber.Map{
		"message": "delete user successful.",
	})
}
