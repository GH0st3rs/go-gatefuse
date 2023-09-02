package rest

import (
	"go-gatefuse/src/config"
	"go-gatefuse/src/storage"

	"github.com/gofiber/fiber/v2"
)

func UiInit(app *fiber.App) {
	app.Get("/", loginRequired, func(c *fiber.Ctx) error { return c.Redirect("/dashboard") })
	app.Get("/dashboard", loginRequired, DashboardHandler)
	app.Get("/settings", loginRequired, GetSettingsHandler)
	app.Post("/settings", loginRequired, PostSettingsHandler)
}

func DashboardHandler(c *fiber.Ctx) error {
	return c.Render("home/dashboard", fiber.Map{
		"settings": config.Settings,
	})
}

func GetSettingsHandler(c *fiber.Ctx) error {
	return c.Render("home/settings", fiber.Map{
		"settings": config.Settings,
	})
}

func PostSettingsHandler(c *fiber.Ctx) error {
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
}
