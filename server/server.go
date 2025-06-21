// server handles the transport layer for communicating with ron-mod-manager's core logic
// currently the only supported transport is HTTP
// It may be moved to grpc when i get around to it, maybe in the next 10 years :D
package server

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"
)

type ServerConf struct {
	Addr string
}

// CreateHTTPServer starts the HTTP server supplying an instance of the application db client and core logic handler
func CreateHTTPServer(db *ent.Client, h handler.Handler, conf ServerConf) error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ron-mm server running")
	})

	app.Post("/actions/add", func(c *fiber.Ctx) error {
		archivePath := c.Query("archivePath")
		name := c.Query("name")

		if archivePath == "" {
			return c.Status(422).SendString("archivePath must be provided")
		}

		if name == "" {
			return c.Status(422).SendString("name must be provided")
		}

		archive, err := h.AddMod(archivePath, name)
		if err != nil {
			slog.Error("Error adding mod")
			slog.Error(err.Error())
			return err
		}

		return c.JSON(archive)
	})

	if err := app.Listen(conf.Addr); err != nil {
		return err
	}

	return nil
}
