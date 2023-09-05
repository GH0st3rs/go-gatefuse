package rest

import (
	"go-gatefuse/src/config"
	"go-gatefuse/src/storage"

	"github.com/gofiber/fiber/v2"
)

func UiInit(app *fiber.App) {
	app.Get("/", loginRequired, func(c *fiber.Ctx) error { return c.Redirect("/dashboard") })
	app.Get("/dashboard", loginRequired, DashboardHandler)
	app.Get("/settings", loginRequired, SettingsHandler)
	app.Post("/settings/save_credentials", loginRequired, SaveCredsHandler)
	app.Post("/settings/save_settings", loginRequired, SaveSettingsHandler)
}

func DashboardHandler(c *fiber.Ctx) error {
	return c.Render("home/dashboard", fiber.Map{
		"settings": config.Settings,
	})
}

func SettingsHandler(c *fiber.Ctx) error {
	return c.Render("home/settings", fiber.Map{
		"settings": config.Settings,
	})
}

func SaveSettingsHandler(c *fiber.Ctx) error {
	var request config.AppSettings
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	// Update settings based on request
	config.Settings.MainDomain = request.MainDomain
	config.Settings.NginxConfPath = request.NginxConfPath
	config.Settings.UnboundConfPath = request.UnboundConfPath
	config.Settings.UnboundRemote = request.UnboundRemote
	config.Settings.UnboundRemoteHost = request.UnboundRemoteHost
	//Save settings
	if err := storage.SaveAppSettings(config.SqliteStorage.Conn()); err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	return c.Redirect("/settings")
}

func SaveCredsHandler(c *fiber.Ctx) error {
	var request config.AppSettings
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	// Update settings based on request
	config.Settings.Username = request.Username
	config.Settings.Password = request.Password
	//Save settings
	if err := storage.SaveAppSettings(config.SqliteStorage.Conn()); err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	return c.Redirect("/settings")
}
