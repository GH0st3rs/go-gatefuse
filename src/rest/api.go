package rest

import (
	"fmt"
	"go-gatefuse/src/config"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

func generateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	letterLength := len(letterRunes)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(letterLength)]
	}
	return string(b)
}
func RestApiInit(app *fiber.App) {
	app.Get("/generate_domain", func(c *fiber.Ctx) error {
		domain := generateRandomString(8)
		return c.JSON(fiber.Map{
			"status":   true,
			"response": fmt.Sprintf("%s.%s", domain, config.Settings.MainDomain),
		})
	})
}
