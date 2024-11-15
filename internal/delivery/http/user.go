package http

import (
	"net/http"
	"time"

	"github.com/Erwanph/be-wan-central-lab/internal/model"
	"github.com/Erwanph/be-wan-central-lab/internal/usecase"
	"github.com/Erwanph/be-wan-central-lab/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
	Config  *viper.Viper
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger, config *viper.Viper) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
		Config:  config,
	}
}

func (c *UserController) RegisterUser(ctx *fiber.Ctx) error {
	request := new(model.RegisterRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed parse register request")
		err = util.ErrInternalDefault
		ctx.Status(err.(util.CustomError).StatusCode())
		return ctx.JSON(model.NewWebResponse("Failed to register a user", err, nil))
	}
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		if customErr, ok := err.(util.CustomError); ok {
			ctx.Status(customErr.StatusCode())
		} else {
			ctx.Status(fiber.StatusInternalServerError)
		}
		return ctx.JSON(model.NewWebResponse("Failed to register a user", err, nil))
	}
	return ctx.JSON(model.NewWebResponse("Successfully register a user", nil, response))
}
func (c *UserController) VerifyEmailRegister(ctx *fiber.Ctx) error {
	request := new(model.RequestVerifyEmailUsingOtp)
	request.Email = ctx.Query("email")
	request.InputOTP = ctx.Query("otp")
	if request.Email == "" || request.InputOTP == "" {
		return ctx.JSON(model.NewWebResponse("Email and OTP are required", nil, nil))
	}

	res, err := c.UseCase.VerifyOTPRegister(ctx.UserContext(), request)
	if err != nil {
		return ctx.JSON(model.NewWebResponse("Failed to verify OTP", err, nil))
	}

	return ctx.JSON(model.NewWebResponse("Email verified successfully", nil, res))
}
func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed parse login request")
		err = util.ErrInternalDefault
		ctx.Status(err.(util.CustomError).StatusCode())
		return ctx.JSON(model.NewWebResponse("Failed to login", err, nil))
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		ctx.Status(err.(util.CustomError).StatusCode())
		return ctx.JSON(model.NewWebResponse("Failed to login", err, nil))
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	tokenExpiresIn, _ := time.ParseDuration(c.Config.GetString("JWT_EXPIRES_IN"))

	claims["sub"] = response.Email
	claims["exp"] = now.Add(tokenExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(c.Config.GetString("JWT_SECRET")))
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed generating JWT token")
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(model.NewWebResponse("Failed to login", util.ErrInternalDefault, nil))
	}

	response.Token = tokenString

	return ctx.JSON(model.NewWebResponse("Login success", nil, response))
}
func (c *UserController) GetProfiles(ctx *fiber.Ctx) error {
	request := new(model.RequestGetProfiles)
	request.UserEmail = ctx.Locals("user").(string)
	profile, err := c.UseCase.GetProfiles(ctx.Context(), request)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogError: err,
		}).Error("Failed to get user")
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(model.NewWebResponse("Failed to get profiles", err, nil))
	}

	return ctx.JSON(model.NewWebResponse("Succes getting profiles", nil, profile))

}
func (c *UserController) UpdateProfiles(ctx *fiber.Ctx) error {
	request := new(model.RequestUpdateProfile)
	request.UserEmail = ctx.Locals("user").(string)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed to parse update user request")
		err = util.ErrInternalDefault
		ctx.Status(err.(util.CustomError).StatusCode())
		return ctx.JSON(model.NewWebResponse("Failed to parse user request", err, nil))
	}

	if request.NewName == "" && request.NewPassword == "" && request.NewEmail == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one field must be provided to update",
		})
	}
	if request.OldPassword == "" && request.NewPassword != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Old password and new password required if wanna update password on profiles",
		})
	}
	if request.NewEmail != "" && request.OldPassword == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password required if wanna update email on profiles",
		})
	}

	updatedUser, err := c.UseCase.UpdateProfiles(ctx.Context(), request)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogError: err,
		}).Error("Failed to update user in usecase")
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(model.NewWebResponse("Failed to update User", err, nil))
	}
	return ctx.JSON(model.NewWebResponse("Profiles has been updated", nil, updatedUser))
}
