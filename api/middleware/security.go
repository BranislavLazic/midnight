package middleware

import (
	sess "github.com/branislavlazic/midnight/api/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"net/http"
)

func Authenticated(sessionStore *session.Store, cookieSecretKey string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		secureCookie := ctx.Cookies(sess.SecureCookieName)
		ok := sess.VerifySessionID(secureCookie, cookieSecretKey)
		if !ok {
			return ctx.SendStatus(http.StatusUnauthorized)
		}
		s, err := sessionStore.Get(ctx)
		if err != nil {
			return ctx.SendStatus(http.StatusUnauthorized)
		}
		data := s.Get(sess.SecureSessionStoreKey)
		if data == nil {
			return ctx.SendStatus(http.StatusUnauthorized)
		}
		return ctx.Next()
	}
}
