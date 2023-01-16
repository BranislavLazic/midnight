package middleware

import (
	sess "github.com/branislavlazic/midnight/api/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"net/http"
)

type Authenticator struct {
	sessionStore    *session.Store
	cookieSecretKey string
}

func NewAuthenticator(sessionStore *session.Store, cookieSecretKey string) *Authenticator {
	return &Authenticator{sessionStore: sessionStore, cookieSecretKey: cookieSecretKey}
}

func (a *Authenticator) Authenticated(next func(ctx *fiber.Ctx) error) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		secureCookie := ctx.Cookies(sess.SecureCookieName)
		ok := sess.VerifySessionID(secureCookie, a.cookieSecretKey)
		if !ok {
			return ctx.SendStatus(http.StatusUnauthorized)
		}
		s, err := a.sessionStore.Get(ctx)
		if err != nil {
			return ctx.SendStatus(http.StatusUnauthorized)
		}
		data := s.Get(sess.SecureSessionStoreKey)
		if data == nil {
			return ctx.SendStatus(http.StatusUnauthorized)
		}
		return next(ctx)
	}
}
