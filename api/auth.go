package api

import (
	secureSession "github.com/branislavlazic/midnight/api/session"
	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthRoutes struct {
	userRepo     model.UserRepository
	sessionStore *session.Store
}

func NewAuthRoutes(userRepo model.UserRepository, sessionStore *session.Store) *AuthRoutes {
	return &AuthRoutes{userRepo: userRepo, sessionStore: sessionStore}
}

// Login godoc
// @Summary Login
// @Param loginRequest body model.LoginRequest true "Login request body"
// @Failure 400,401,404,422
// @Success 200
// @Router /v1/login [post]
func (ar *AuthRoutes) Login(ctx *fiber.Ctx) error {
	var loginRequest *model.LoginRequest
	if err := ctx.BodyParser(&loginRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as login request")
		return ctx.SendStatus(http.StatusBadRequest)
	}
	err := validator.New().Struct(loginRequest)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	user, err := ar.userRepo.GetByEmail(loginRequest.Email)
	if err != nil {
		log.Debug().Err(err).Msgf("failed to find the user by email %s", loginRequest.Email)
		return ctx.SendStatus(http.StatusNotFound)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		log.Debug().Err(err).Msgf("incorrect password")
		return ctx.SendStatus(http.StatusUnauthorized)
	}
	sess, err := ar.sessionStore.Get(ctx)
	if err != nil {
		return ctx.SendStatus(http.StatusUnauthorized)
	}
	sess.Set(secureSession.SecureSessionStoreKey, user.Email)
	if err := sess.Save(); err != nil {
		log.Error().Err(err).Msgf("failed to set the session")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.Status(http.StatusOK).JSON(user)
}

// Logout godoc
// @Summary Logout
// @Failure 500
// @Success 200
// @Router /v1/logout [post]
func (ar *AuthRoutes) Logout(ctx *fiber.Ctx) error {
	err := ar.sessionStore.Reset()
	if err != nil {
		log.Error().Err(err).Msgf("failed to remove the session")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	expireCookie(ctx)
	return ctx.SendStatus(http.StatusOK)
}

func expireCookie(ctx *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = secureSession.SecureCookieName
	cookie.Expires = time.Now().Add(-3 * time.Second)
	ctx.Cookie(cookie)
}
