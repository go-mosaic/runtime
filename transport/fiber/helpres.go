package fiber

import (
	"github.com/go-mosaic/runtime/transport"
	"github.com/gofiber/fiber/v3"
)

func convertSameSite(sameSite transport.SameSite) string {
	switch sameSite {
	case transport.SameSiteLaxMode:
		return fiber.CookieSameSiteLaxMode
	case transport.SameSiteNoneMode:
		return fiber.CookieSameSiteNoneMode
	case transport.SameSiteStrictMode:
		return fiber.CookieSameSiteStrictMode
	default:
		return fiber.CookieSameSiteDisabled
	}
}
