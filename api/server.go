package api

import (
	db "smutaxi/db/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	store  *db.Store
	router *fiber.App
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := fiber.New()

	router.Use(logger.New())
	router.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON("pong")
	})

	// User router
	router.Get("/users", server.getUsers)
	router.Get("/users/:id", server.getUser)
	router.Post("/users", server.createUser)
	router.Delete("/users/:id", server.deleteUser)

	// Room router
	router.Get("/rooms", server.getAllRooms)
	router.Get("/rooms/:id", server.getRoom)
	router.Patch("/rooms/:id", server.addMember)
	router.Post("/rooms/:roomID", server.EnterRoom)
	router.Post("/rooms", server.createRoom)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}

func errorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"error": err.Error(),
	}
}
