package server

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
)

func CreateServer(db *ent.Client) (*fiber.App, error) {
	app := fiber.New()
	appConf, err := appconfig.Read()
	if err != nil {
		return nil, err
	}
	iohandler := handlerio.NewFileSystemHandler()
	h := handler.NewHandler(db, appConf, iohandler)

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

	if err := app.Listen(":5000"); err != nil {
		return nil, err
	}

	return app, nil
}
