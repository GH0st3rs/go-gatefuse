package rest

import (
	"go-gatefuse/src/config"
	"go-gatefuse/src/storage"

	"github.com/gofiber/fiber/v2"
)

func UiInit(app *fiber.App) {
	app.Get("/index", func(c *fiber.Ctx) error {
		return c.Render("home/index", fiber.Map{})
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.Render("home/dashboard", fiber.Map{
			"settings": config.Settings,
		})
	})

	app.Get("/settings", func(c *fiber.Ctx) error {
		return c.Render("home/settings", fiber.Map{
			"settings": config.Settings,
		})
	})

	app.Post("/settings", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&config.Settings); err != nil {
			return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
		}
		if err := storage.SaveAppSettings(config.SqliteStorage.Conn()); err != nil {
			return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
		}
		return c.Redirect("/settings")
	})

}
