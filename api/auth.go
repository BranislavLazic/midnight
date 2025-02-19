package api

import (
	"net/http"
	"time"

	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/config"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const tokenValidityDuration = time.Minute * 30

type jwtCustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type jwtAuthResponse struct {
	AuthUser    model.User `json:"authUser"`
	AccessToken string     `json:"accessToken"`
}

type AuthRoutes struct {
	repo   *postgres.Repository
	config *config.AppConfig
}

func NewAuthRoutes(repo *postgres.Repository, config *config.AppConfig) *AuthRoutes {
	return &AuthRoutes{repo: repo, config: config}
}

// Login godoc
// @Summary Login
// @Param loginRequest body model.LoginRequest true "Login request body"
// @Failure 400,401,404,422
// @Success 200
// @Router /v1/login [post]
func (ar *AuthRoutes) Login(ctx echo.Context) error {
	var loginRequest *model.LoginRequest
	if err := ctx.Bind(&loginRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as login request")
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := validator.New().Struct(loginRequest)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	user, err := ar.repo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		log.Debug().Err(err).Msgf("failed to find the user by email %s", loginRequest.Email)
		return ctx.NoContent(http.StatusNotFound)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		log.Debug().Err(err).Msgf("incorrect password")
		return ctx.NoContent(http.StatusUnauthorized)
	}
	claims := &jwtCustomClaims{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenValidityDuration)),
		},
	}
	jwtTok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtTok.SignedString([]byte(ar.config.SessionSecret))
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}
	return ctx.JSON(http.StatusOK, jwtAuthResponse{
		AuthUser:    *user,
		AccessToken: token,
	})
}
