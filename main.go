package main

import (
	"errors"
	"flag"
	"fmt"
	"go-gatefuse/src/config"
	"go-gatefuse/src/rest"
	"go-gatefuse/src/storage"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/django/v3"
)

type Route struct {
	Active      bool
	Target      string
	Source      string
	Description string
	UUID        string
}

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	// Retrieve the custom status code if it's a *fiber.Error
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	// Check for status, if it's a 404, render the 404 template.
	if code == fiber.StatusNotFound {
		return ctx.Status(code).Render("home/page-404", fiber.Map{})
	} else if code == fiber.StatusInternalServerError {
		return ctx.Status(code).Render("home/page-500", fiber.Map{})
	}
	return fiber.DefaultErrorHandler(ctx, err)
}

func main() {
	flag.Parse()
	// Create a new engine
	engine := django.New("./templates", ".html")
	engine.Reload(true)
	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
		// Override default error handler
		ErrorHandler: CustomErrorHandler,
	})
	app.Static("/static/assets", "./static/assets")
	// Enable debug output
	if *config.AppDebug {
		app.Use(logger.New())
	}
	// Associate favicon image
	app.Use(favicon.New(favicon.Config{
		File: "./static/assets/images/logo.png",
		URL:  "/favicon.ico",
	}))
	// Add CSRF cookie token
	app.Use(csrf.New(csrf.Config{
		Storage:    config.SqliteStorage,
		ContextKey: "csrf_token",
		CookieName: "csrf_token",
		KeyLookup:  "cookie:csrf_token",
	}))
	// Use cache middleware with a global expiration time of 10 minutes
	if *config.UseCache {
		// Next -  defines a function to skip the middleware.
		app.Use(cache.New(cache.Config{Next: nil, Expiration: 10 * time.Minute}))
	}
	// Initialize database
	if *config.AppInit {
		if err := storage.InitializeDatabaseTables(config.SqliteStorage.Conn()); err != nil {
			log.Fatalln("Failed to create additional tables: ", err.Error())
		}
		return
	}

	if err := storage.LoadAppSettings(config.SqliteStorage.Conn(), &config.Settings); err != nil {
		log.Fatalln("Failed to load the settings: ", err.Error())
	}

	// Register new handlers
	rest.AuthInit(app)
	rest.UiInit(app)
	rest.RestApiInit(app)
	rest.GateInit(app)

	// Start server listener
	url := fmt.Sprintf("%s:%d", *config.AppHost, *config.AppPort)
	if err := app.Listen(url); err != nil {
		log.Fatalln("Failed to start: ", err.Error())
	}
}
