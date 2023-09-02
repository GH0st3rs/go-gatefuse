package rest

import (
	"go-gatefuse/src/config"
	"go-gatefuse/src/nginx"
	"go-gatefuse/src/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GateInit(app *fiber.App) {
	app.Get("/list", ListHandler)
	app.Post("/create", CreateHandler)
	app.Post("/update", UpdateHandler)
	app.Post("/toggle", ToggleHandler)
	app.Post("/delete", DeleteHandler)
}

func ListHandler(c *fiber.Ctx) error {
	items, err := storage.RetrieveAllGateRecords(config.SqliteStorage.Conn())
	if err != nil {
		return c.JSON(fiber.Map{"status": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": true, "routes": items})
}

func CreateHandler(c *fiber.Ctx) error {
	var request config.GateRecord
	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not parse request: " + err.Error(),
		})
	}

	request.UUID = uuid.NewString()
	// Add new record to Database
	if err := storage.AddNewRecord(config.SqliteStorage.Conn(), request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not add the record to database: " + err.Error(),
		})
	}
	// Create a configuration file if activated
	if err := nginx.SaveNginxConfig(request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not create configuration: " + err.Error(),
		})
	}
	// Reload Nginx
	if err := nginx.ReloadNginx(); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not reload nginx: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"status": true, "response": request})
}

func UpdateHandler(c *fiber.Ctx) error {
	var request config.GateRecord
	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not parse request: " + err.Error(),
		})
	}
	// Update record in database
	if err := storage.UpdateGateRecord(config.SqliteStorage.Conn(), request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not update the record: " + err.Error(),
		})
	}
	// Create a configuration file if activated
	if err := nginx.SaveNginxConfig(request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not create configuration: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"status": true, "response": request})
}

func ToggleHandler(c *fiber.Ctx) error {
	var request config.GateRecord
	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not parse request: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"status": true, "response": request})
}

func DeleteHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": true})
}
