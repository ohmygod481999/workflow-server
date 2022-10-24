package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RequireUser(c *fiber.Ctx) error {
	if user := c.Locals("user"); user == nil {
		return c.Status(403).SendString("Forbidden")
	}
	return c.Next()
}

func RequireCallBotAgent(c *fiber.Ctx) error {
	if agent := c.Locals("agent"); agent == nil {
		return c.Status(403).SendString("Forbidden")
	}
	return c.Next()
}

func RequireUserOrCallbotAgent(c *fiber.Ctx) error {
	if user := c.Locals("user"); user != nil {
		return c.Next()
	} else if agent := c.Locals("agent"); agent != nil {
		return c.Next()
	} else {
		return c.Status(403).SendString("Forbidden")
	}
}
