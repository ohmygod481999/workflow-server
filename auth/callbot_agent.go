package auth

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type Deployment struct {
	Id     *uuid.UUID `json:"id"`
	UserId string     `json:"user_id"`
	Status string     `json:"status"`
}

type ResolvedApiKey struct {
	Deployment Deployment `json:"deployment"`
}

func NewCallbotAgentAuthenticator(cfg Config) fiber.Handler {
	cache := cache.New(5*time.Minute, 10*time.Minute)

	return func(c *fiber.Ctx) error {
		url := fmt.Sprintf("%s/api-key/inspect", cfg.VintalkServicesUri)
		apiKey := c.Get("X-API-KEY")

		var agent ResolvedApiKey

		if apiKey == "" {
			c.Locals("agent", nil)
			return c.Next()
		}
		if u, found := cache.Get(apiKey); found {
			agent = u.(ResolvedApiKey)
		} else {
			res, err := client.R().
				SetQueryParam("api_key", apiKey).
				SetResult(&agent).
				Get(url)
			if err != nil {
				return err
			}
			if res.StatusCode() != 200 || agent.Deployment.Id == nil {
				c.Locals("agent", nil)
				return c.Next()
			} else {
				cache.Set(apiKey, agent, 5*time.Minute)
			}
		}
		c.Locals("agent", &agent)
		return c.Next()
	}
}
