package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var client = resty.New()
var pFalse = false

func getBearerToken(h string) string {
	parts := strings.Split(h, " ")
	if parts[0] != "Bearer" {
		return ""
	}
	n := len(parts)
	return parts[n-1]
}

type User struct {
	Id          string `json:"id"`
	DisplayName string `json:"display_name"`
	IsBlocked   *bool  `json:"is_blocked"`
	Email       string `json:"email"`
	Token       string
}
type PassportResponse struct {
	Data User `json:"data"`
}

func readRootTokens() map[string]bool {
	env := viper.GetString("root_tokens")
	fmt.Println(env)
	env = strings.ReplaceAll(env, "'", "\"")
	var rootTokens []string
	if err := json.Unmarshal([]byte(env), &rootTokens); err != nil {
		panic(err)
	}
	result := make(map[string]bool)
	for _, t := range rootTokens {
		result[t] = true
	}
	return result
}

func NewPassportAuthenticator(cfg Config) fiber.Handler {
	cache := cache.New(5*time.Minute, 10*time.Minute)

	rootTokens := readRootTokens()

	return func(c *fiber.Ctx) error {
		var token string
		var pp PassportResponse
		switch {
		case c.Get("authorization") != "":
			token = getBearerToken(c.Get("authorization"))
		case c.Get("token") != "":
			token = c.Get("token")
		case c.Query("token") != "":
			token = c.Query("token")
		case c.Cookies("token") != "":
			token = c.Cookies("token")
		}
		if token == "" {
			c.Locals("user", nil)
			return c.Next()
		}

		if _, found := rootTokens[token]; found {
			user := &User{
				Id:          "root",
				DisplayName: "root",
				IsBlocked:   &pFalse,
				Email:       "",
				Token:       token,
			}
			c.Locals("user", user)
			return c.Next()
		}

		if u, found := cache.Get(token); found {
			pp = u.(PassportResponse)
		} else {
			url := fmt.Sprintf("%s/users/user", cfg.PassportUri)

			_, err := client.R().
				SetHeader("cookie", fmt.Sprintf("token=%s", token)).
				SetResult(&pp).
				Get(url)

			if err != nil {
				return err
			}
		}
		user := &pp.Data
		user.Token = token

		if user.Id == "" {
			c.Locals("user", nil)
		} else {
			cache.Set(token, pp, 5*time.Minute)
			c.Locals("user", user)
		}
		return c.Next()
	}
}
