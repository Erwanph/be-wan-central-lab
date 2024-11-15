package converter

import (
	"github.com/Erwanph/be-wan-central-lab/internal/entity"
	"github.com/Erwanph/be-wan-central-lab/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewRegisterResponse(user *entity.User) *model.RegisterResponse {
	return &model.RegisterResponse{
		Name:            user.Name,
		Email:           user.Email,
		CreatedAt:       primitive.NewDateTimeFromTime(user.CreatedAt),
		IsEmailVerified: user.IsEmailVerified,
	}
}
func NewResponseVerifyEmailUsingOtp(user *entity.User) *model.ResponseVerifyEmailUsingOtp {
	return &model.ResponseVerifyEmailUsingOtp{
		Name:      user.Name,
		Email:     user.Email,
		UpdatedAt: primitive.NewDateTimeFromTime(*user.UpdatedAt),
	}
}
func NewVerifyAuthResponse(user *entity.User) *model.VerifyAuthResponse {
	return &model.VerifyAuthResponse{
		Email: user.Email,
	}
}

func NewLoginResponse(user *entity.User) *model.LoginResponse {
	return &model.LoginResponse{
		Email: user.Email,
	}
}

func NewUpdateUserResponse(user *entity.User) *model.ResponseUpdateProfile {
	return &model.ResponseUpdateProfile{
		Name:        user.Name,
		Email:       user.Email,
		NewPassword: user.Password,
		UpdatedAt:   user.UpdatedAt,
	}
}
func NewGetUserResponse(user *entity.User) *model.ResponseGetProfiles {
	return &model.ResponseGetProfiles{
		Name:      user.Name,
		Email:     user.Email,
		Score: 	user.Score,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

