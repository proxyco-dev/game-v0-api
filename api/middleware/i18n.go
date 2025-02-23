package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	DefaultLanguage = "ka"
)

func I18nMiddleware(bundle *i18n.Bundle) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accept := c.Get("Accept-Language")
		if accept == "" {
			accept = DefaultLanguage
		}

		localizer := i18n.NewLocalizer(bundle, accept)
		c.Locals("localizer", localizer)

		return c.Next()
	}
}
