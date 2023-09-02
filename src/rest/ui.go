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
		var request config.AppSettings
		if err := c.BodyParser(&request); err != nil {
			return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
		}
		// Update settings based on request
		switch request.RequestType {
		case "appCredentialsForm":
			config.Settings.Username = request.Username
			config.Settings.Password = request.Password
		case "appSettingsForm":
			config.Settings.MainDomain = request.MainDomain
			config.Settings.NginxConfPath = request.NginxConfPath
		}
		//Save settings
		if err := storage.SaveAppSettings(config.SqliteStorage.Conn()); err != nil {
			return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
		}
		return c.Redirect("/settings")
	})

	app.Get("/home/:page", func(c *fiber.Ctx) error {
		page := c.Params("page")
		return c.Render("home/"+page, fiber.Map{})
	})
}
