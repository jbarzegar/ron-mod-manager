// server handles the transport layer for communicating with ron-mod-manager's core logic
// currently the only supported transport is HTTP
// It may be moved to grpc when i get around to it, maybe in the next 10 years :D
package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"

	"github.com/jbarzegar/ron-mod-manager/internal/actions"
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

	app.Get("/mods", func(c *fiber.Ctx) error {
		staged, err := h.GetAllMods()
		if err != nil {
			return err
		}

		return c.JSON(staged)
	})

	app.Get("/archives", func(c *fiber.Ctx) error {
		var req *actions.GetArchiveRequest

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		fmt.Println(req.Untracked)
		archives, err := h.GetArchives(req)

		if err != nil {
			return err
		}

		return c.JSON(archives)
	})

	app.Get("/staged", func(c *fiber.Ctx) error {
		staged, err := h.GetStagedMods()
		if err != nil {
			return err
		}

		return c.JSON(staged)
	})

	// List the staging enviroment see what mods are active
	// app.Get("/staging")

	app.Post("/actions/add", func(c *fiber.Ctx) error {
		// archivePath := c.Query("archivePath")
		// name := c.Query("name")

		var body *actions.AddRequest
		err := json.Unmarshal(c.Body(), &body)
		if err != nil {
			return err
		}

		if body.ArchivePath == "" {
			return c.Status(http.StatusUnprocessableEntity).
				SendString("archivePath must be provided")
		}

		if body.Name == "" {
			return c.Status(http.StatusUnprocessableEntity).
				SendString("name must be provided")
		}

		archive, err := h.AddArchive(body.ArchivePath, body.Name)
		if err != nil {
			slog.Error("Error adding mod")
			slog.Error(err.Error())
			return err
		}

		return c.Status(http.StatusOK).JSON(archive)
	})

	app.Post("/actions/install", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var body actions.InstallRequest
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			return err
		}

		mUid, err := uuid.Parse(body.ModVersion)
		if err != nil {
			return err
		}

		if err := h.InstallMod(body.ModID, mUid, body.Choices); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)

	})

	app.Delete("/actions/uninstall", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var body actions.UninstallModRequest
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			return err
		}

		if err := h.UninstallMod(body.ModIds); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

	app.Delete("/actions/delete", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var body actions.DeleteModRequest
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			return err
		}

		if err := h.DeleteMod(body); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	})

	if err := app.Listen(conf.Addr); err != nil {
		return err
	}

	return nil
}
