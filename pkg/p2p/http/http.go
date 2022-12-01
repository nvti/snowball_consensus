package http

import (
	"snowball/models"

	"github.com/gofiber/fiber/v2"
)

type peerFoundHandler func(peerAddress string)

func CreateHttpServer(host string, port int, handler peerFoundHandler) *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		req := &models.NewNodeHook{}
		if err := c.BodyParser(req); err != nil {
			return err
		}
		handler(req.Address)

		return c.SendStatus(fiber.StatusOK)
	})

	return app
}
