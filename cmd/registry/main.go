package main

import (
	"flag"
	"snowball/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var peers = []string{}
var host string
var port int

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Listen on interface")
	flag.IntVar(&port, "host", 5001, "Listen on port")

	flag.Parse()
}

func main() {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		req := &models.RegisterNodeReq{}
		if err := c.BodyParser(req); err != nil {
			return err
		}

		nodeAddress := c.IP() + ":" + strconv.Itoa(req.Port)
		peers = append(peers, nodeAddress)

		return c.JSON(peers)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(peers)
	})

	app.Listen(host + ":" + strconv.Itoa(port))
}
