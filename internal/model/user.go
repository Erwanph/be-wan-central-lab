package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type RegisterRequestCustomRoles struct {
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RolesID  string `json:"roles_id" validate:"required"`
}
type RegisterResponse struct {
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	CreatedAt       primitive.DateTime `json:"created_at"`
	IsEmailVerified bool               `json:"is_email_verified"`
}
type RequestVerifyEmailUsingOtp struct {
	Email    string `json:"email" validate:"required,email"`
	InputOTP string `json:"otp" validate:"required"`
}
type ResponseVerifyEmailUsingOtp struct {
	Email     string             `json:"email" validate:"required,email"`
	Name      string             `json:"name"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type VerifyAuthRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyAuthResponse struct {
	Email string `json:"email"`
}

type RequestOTPResetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestUpdateProfile struct {
	UserEmail   string `json:"email" validate:"required,email"`
	NewName     string `json:"new_name"`
	NewEmail    string `json:"new_email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ResponseUpdateProfile struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	NewPassword string     `json:"new_password"`
	NewJWTToken string     `json:"token"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type RequestGetProfiles struct {
	UserEmail string `json:"email" validate:"required,email"`
}

type ResponseGetProfiles struct {
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
