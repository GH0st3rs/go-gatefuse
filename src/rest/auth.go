package rest

import (
	"go-gatefuse/src/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username   string
	Password   string
	RememberMe string
}

func loginRequired(c *fiber.Ctx) error {
	session, err := config.SessionStorage.Get(c)
	if err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	if session.Get("username") != config.Settings.Username {
		return c.Redirect("/login")
	}
	return c.Next()
}

func AuthInit(app *fiber.App) {
	app.Get("/login", GetLoginHandler)
	app.Post("/login", PostLoginHandler)
	app.Get("/logout", GetLogoutHandler)
}

func GetLoginHandler(c *fiber.Ctx) error {
	session, err := config.SessionStorage.Get(c)
	if err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	if session.Get("username") == config.Settings.Username {
		return c.Redirect("/dashboard")
	}
	return c.Render("accounts/login", fiber.Map{
		"csrf_token": c.Locals("csrf_token"),
	})
}

func PostLoginHandler(c *fiber.Ctx) error {
	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}

	session, err := config.SessionStorage.Get(c)
	if err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	if request.Username == "admin" && request.Password == "admin" {
		session.Set("username", request.Username)
		if request.RememberMe == "on" {
			session.SetExpiry(7 * 24 * time.Hour)
		}
		// Save session
		if err := session.Save(); err != nil {
			return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
		}
		// Success login
		return c.Redirect("/dashboard")
	}
	// Wrong login
	return c.Render("accounts/login", fiber.Map{
		"csrf_token": c.Locals("csrf_token"),
		"msg":        "Wrong user or password",
	})
}

func GetLogoutHandler(c *fiber.Ctx) error {
	session, err := config.SessionStorage.Get(c)
	if err != nil {
		return c.Status(500).Render("home/page-500", fiber.Map{"msg": err.Error()})
	}
	session.Destroy()
	return c.Redirect("/login")
}
