package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Erwanph/be-wan-central-lab/internal/model"
	"github.com/Erwanph/be-wan-central-lab/internal/usecase"
	"github.com/Erwanph/be-wan-central-lab/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AuthMiddleware struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
	Config  *viper.Viper
}

func NewAuthMiddleware(log *logrus.Logger, useCase *usecase.UserUseCase, config *viper.Viper) *AuthMiddleware {
	return &AuthMiddleware{
		Log:     log,
		UseCase: useCase,
		Config:  config,
	}
}

func (m *AuthMiddleware) CheckSession(ctx *fiber.Ctx) error {
	var tokenString string
	authorization := ctx.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if ctx.Cookies("token") != "" {
		tokenString = ctx.Cookies("token")
	}

	if tokenString == "" {
		return ctx.JSON(model.NewWebResponse("Authentication failed", util.ErrNotLoginYet, nil))
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(m.Config.GetString("JWT_SECRET")), nil
	})
	if err != nil {
		m.Log.WithFields(logrus.Fields{
			util.LogError: err,
		}).Error("Failed parsing JWT token")
		ctx.Status(http.StatusUnauthorized)
		return ctx.JSON(model.NewWebResponse("Authentication failed", util.ErrInternalDefault, nil))
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return ctx.JSON(model.NewWebResponse("Authentication failed", util.ErrInvalidToken, nil))
	}

	request := &model.VerifyAuthRequest{
		Email: claims["sub"].(string),
	}
	user, err := m.UseCase.Verify(ctx.UserContext(), request)

	if errors.Is(err, util.ErrInvalidCredential) {
		return ctx.JSON(model.NewWebResponse("Authentication failed", util.ErrInvalidToken, nil))
	}

	ctx.Locals("user", user.Email)

	return ctx.Next()
}

