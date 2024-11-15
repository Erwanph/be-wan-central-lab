package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Erwanph/be-wan-central-lab/internal/entity"
	"github.com/Erwanph/be-wan-central-lab/internal/model"
	"github.com/Erwanph/be-wan-central-lab/internal/model/converter"
	"github.com/Erwanph/be-wan-central-lab/internal/repository"
	"github.com/Erwanph/be-wan-central-lab/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	Config         *viper.Viper
}

func NewUserUseCase(logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, config *viper.Viper) *UserUseCase {
	return &UserUseCase{
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		Config:         config,
	}
}
func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	if err := util.ValidateRequestRegister(request); err != nil {
		return nil, err
	}
	request.Name = util.GetDefaultName(request.Name)
	total, err := c.UserRepository.CountByEmail(ctx, request.Email)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed to count user by Email in database")
		return nil, util.ErrInternalDefault
	}
	if total > 0 {
		return nil, util.ErrUserAlreadyExist
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed to hash password")
		return nil, err
	}
	user := &entity.User{
		Name:            request.Name,
		Email:           strings.ToLower(request.Email),
		Password:        string(password),
		CreatedAt:       util.NowInWIB(),
		IsEmailVerified: true, //set untuk otp
	}
	user.SecretKey, err = util.GenerateSecretKey(user.Email)
	if err != nil {
		return nil, err
	}

	OTP, err := util.GenerateOTPRegister()
	if err != nil {
		return nil, err
	}

	user.OTP, err = util.HashOTP(OTP)
	if err != nil {
		return nil, err
	}

	// err = util.SendOTPRegister(user.Email, OTP, c.Config)
	// if err != nil {
	// 	return nil, err
	// }

	user.OTPExpiresAt = util.NowInWIB().Add(15 * time.Minute)
	if err = c.UserRepository.CreateDefaultUser(ctx, user); err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed create user to database")
		return nil, err
	}

	user.Score = 0

	return converter.NewRegisterResponse(user), nil
}

func (c *UserUseCase) UpdateScore(ctx context.Context, email string, score int) (*model.UpdateScoreResponse, error) {
	if !util.IsValidEmail(email) {
		return nil, util.ErrInvalidEmail
	}
	user, err := c.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, util.ErrInternalDefault
	}
	if user.Score > 0 {
		return nil, errors.New("assignment already worked")
	}
	user.Score += score // Increment existing score

	err = c.UserRepository.UpdateScore(ctx, user)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"email": email,
			"score": score,
			"error": err,
		}).Error("Failed to update user score")
		return nil, err
	}

	return &model.UpdateScoreResponse{
		Email: user.Email,
		Score: user.Score,
	}, nil
}

func (c *UserUseCase) VerifyOTPRegister(ctx context.Context, request *model.RequestVerifyEmailUsingOtp) (*model.ResponseVerifyEmailUsingOtp, error) {
	user, err := c.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil || user == nil {
		return nil, err
	}
	isOTPValid := util.VerifyOTPRegister(user.OTP, request.InputOTP)
	isOTPExpired := util.NowInWIB().After(user.OTPExpiresAt)
	if !isOTPValid && isOTPExpired {
		return nil, errors.New("invalid OTP and OTP has expired")
	} else if !isOTPValid {
		return nil, errors.New("invalid OTP")
	} else if isOTPExpired {
		return nil, errors.New("OTP has expired")
	}

	user.IsEmailVerified = true
	user.OTP = ""
	user.OTPExpiresAt = time.Time{}

	if err = c.UserRepository.VerifiedEmailUser(ctx, user); err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("failed verifed email register")
		return nil, err
	}
	updated_time := util.NowInWIB()
	user.UpdatedAt = &updated_time
	return converter.NewResponseVerifyEmailUsingOtp(user), nil

}
func (c *UserUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	err := c.Validate.Struct(request)
	if err != nil {
		return nil, util.NewCustomError(err)
	}
	if err := util.ValidateRequestLogin(request); err != nil {
		return nil, err
	}

	user, err := c.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed to find user by Email in database")
		return nil, err
	}

	if user == nil {
		return nil, util.ErrInvalidCredential
	}

	if !user.IsEmailVerified {
		return nil, util.NewCustomError(errors.New("email not verified"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, util.ErrInvalidCredential
	}
	return converter.NewLoginResponse(user), nil
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyAuthRequest) (*model.VerifyAuthResponse, error) {
	user, err := c.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			util.LogRequest: request,
			util.LogError:   err,
		}).Error("Failed to find user by Email in database")
		return nil, err
	}

	if user == nil {
		return nil, util.ErrInvalidCredential
	}

	return converter.NewVerifyAuthResponse(user), nil
}
func (c *UserUseCase) GetProfiles(ctx context.Context, request *model.RequestGetProfiles) (*model.ResponseGetProfiles, error) {
	err := c.Validate.Struct(request)
	if err != nil {
		return nil, util.NewCustomError(err)
	}
	user, err := c.UserRepository.FindByEmail(ctx, request.UserEmail)
	if err != nil {
		return nil, err
	}
	return converter.NewGetUserResponse(user), err
}
func (c *UserUseCase) UpdateProfiles(ctx context.Context, request *model.RequestUpdateProfile) (*model.ResponseUpdateProfile, error) {
	err := c.Validate.Struct(request)
	new_jwt_token := ""
	if err != nil {
		return nil, util.NewCustomError(err)
	}
	user, err := c.UserRepository.FindByEmail(ctx, request.UserEmail)
	if err != nil {
		return nil, err
	}

	if request.NewPassword != "" && request.OldPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
		if err != nil {
			return nil, util.ErrOldPasswordNotMatched
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.NewPassword))
		if err == nil {
			return nil, util.ErrSameOldAndNewPassword
		}

		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user.Password = string(newHashedPassword)
	}
	if request.NewEmail != "" {
		if !util.IsValidEmail(request.NewEmail) {
			return nil, util.ErrInvalidEmail
		}
		if !util.IsValidDomain(request.NewEmail) {
			return nil, util.ErrInvalidDomain
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
		if err != nil {
			return nil, util.ErrOldPasswordNotMatched
		}
		if user.Email == request.NewEmail {
			return nil, errors.New("old email and new email are same")
		}
		user.Email = request.NewEmail
		tokenByte := jwt.New(jwt.SigningMethodHS256)

		now := time.Now().UTC()
		claims := tokenByte.Claims.(jwt.MapClaims)
		tokenExpiresIn, _ := time.ParseDuration(c.Config.GetString("JWT_EXPIRES_IN"))
		claims["sub"] = user.Email
		claims["exp"] = now.Add(tokenExpiresIn).Unix()
		claims["iat"] = now.Unix()
		claims["nbf"] = now.Unix()
		new_jwt_token, err = tokenByte.SignedString([]byte(c.Config.GetString("JWT_SECRET")))
		if err != nil {
			return nil, err
		}

	}
	if request.NewName != "" {
		if user.Name == request.NewName {
			return nil, errors.New("old name and new name are same")
		}
		user.Name = request.NewName
	}

	err = c.UserRepository.UpdateProfiles(ctx, user, request.UserEmail)
	if err != nil {
		return nil, err
	}
	updatedUser := converter.NewUpdateUserResponse(user)
	updatedUser.NewJWTToken = new_jwt_token
	return updatedUser, nil
}
