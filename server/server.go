package server

import (
	"github.com/gofiber/fiber"
	"github.com/jbarzegar/ron-mod-manager/ent"
)

func CreateServer(db *ent.Client) (*fiber.App, error) {
	app := fiber.New()

	err := app.Listen(":5000")
	if err != nil {
		return nil, err
	}

	return app, nil
}
