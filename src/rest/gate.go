package rest

import (
	"go-gatefuse/src/config"
	"go-gatefuse/src/nginx"
	"go-gatefuse/src/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeleteRequest struct {
	UUID string `json:"uuid" form:"uuid"`
}

type ToggleRequest struct {
	UUID   string `json:"UUID" form:"UUID"`
	Active bool   `json:"Active" form:"Active"`
}

func GateInit(app *fiber.App) {
	app.Get("/list", loginRequired, ListHandler)
	app.Post("/create", loginRequired, CreateHandler)
	app.Post("/update", loginRequired, UpdateHandler)
	app.Post("/toggle", loginRequired, ToggleHandler)
	app.Post("/delete", loginRequired, DeleteHandler)
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
	return saveAndReload(c, request)
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
	return saveAndReload(c, request)
}

func ToggleHandler(c *fiber.Ctx) error {
	var request ToggleRequest
	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{"status": false, "error": "Could not parse request: " + err.Error()})
	}
	// Retrieve from DB to be sure that the record exists
	record, err := storage.RetrieveOneGateRecord(config.SqliteStorage.Conn(), request.UUID)
	if err != nil {
		return c.JSON(fiber.Map{"status": false, "error": "Could not retrieve the record: " + err.Error()})
	}
	record.Active = request.Active
	// Update record in database
	if err := storage.UpdateGateRecord(config.SqliteStorage.Conn(), record); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not update the record: " + err.Error(),
		})
	}
	if record.Active {
		// Create a configuration file if activated
		if err := nginx.SaveNginxConfig(record); err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"error":  "Could not create configuration: " + err.Error(),
			})
		}
	} else {
		// Delete all configuration files
		if err := nginx.DeleteNginxConfig(record); err != nil {
			return c.JSON(fiber.Map{"status": false, "error": "Could not remove conf files: " + err.Error()})
		}
	}
	// Reload Nginx
	if err := nginx.ReloadNginx(); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not reload nginx: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"status": true, "response": record})
}

func DeleteHandler(c *fiber.Ctx) error {
	var request DeleteRequest
	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{"status": false, "error": "Could not parse request: " + err.Error()})
	}
	// Retrieve from DB to be sure that the record exists
	record, err := storage.RetrieveOneGateRecord(config.SqliteStorage.Conn(), request.UUID)
	if err != nil {
		return c.JSON(fiber.Map{"status": false, "error": "Could not retrieve the record: " + err.Error()})
	}
	// Delete all configuration files
	if err := nginx.DeleteNginxConfig(record); err != nil {
		return c.JSON(fiber.Map{"status": false, "error": "Could not remove conf files: " + err.Error()})
	}
	// Delete record from DB
	if err := storage.DeleteGateRecord(config.SqliteStorage.Conn(), record.UUID); err != nil {
		return c.JSON(fiber.Map{"status": false, "error": "Could not delete record from db: " + err.Error()})
	}
	// Reload Nginx
	if err := nginx.ReloadNginx(); err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"error":  "Could not reload nginx: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"status": true, "response": record.UUID})
}

func saveAndReload(c *fiber.Ctx, record config.GateRecord) error {
	// Create a configuration file if activated
	if err := nginx.SaveNginxConfig(record); err != nil {
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
	return c.JSON(fiber.Map{"status": true, "response": record})
}
