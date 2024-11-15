package util

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CustomError struct {
	Code int
	Err  error
}

func NewCustomError(err error) CustomError {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		return CustomError{Code: http.StatusBadRequest, Err: err}
	}

	return CustomError{Code: http.StatusInternalServerError, Err: err}
}

func (e CustomError) Error() string {
	return e.Err.Error()
}

func (e CustomError) Unwrap() error {
	return e.Err
}

func (e CustomError) StatusCode() int {
	return e.Code
}

var (
	// login error
	ErrInternalDefault   = CustomError{http.StatusInternalServerError, errors.New("something went wrong")}
	ErrInvalidCredential = CustomError{http.StatusUnauthorized, errors.New("invalid credential")}

	ErrInvalidCredentialUpdateUser = CustomError{http.StatusUnauthorized, errors.New("invalid id on request update user")}
	ErrNotLoginYet                 = CustomError{http.StatusUnauthorized, errors.New("you are not logged in")}
	ErrInvalidToken                = CustomError{http.StatusUnauthorized, errors.New("invalid or expired token")}
	ErrInvalidFormatRequest        = CustomError{http.StatusBadRequest, errors.New("invalid format request")}

	//register error
	ErrInvalidEmail     = CustomError{http.StatusBadRequest, errors.New("invalid email format request make sure the format is name@domain")}
	ErrInvalidDomain    = CustomError{http.StatusBadRequest, errors.New("domain email was not valid. Please check your domain again")}
	ErrUserAlreadyExist = CustomError{http.StatusConflict, errors.New("user already exists")}

	ErrOldPasswordNotMatched = CustomError{http.StatusBadRequest, errors.New("old password not matched on database")}
	ErrSameOldAndNewPassword = CustomError{http.StatusBadRequest, errors.New("old password and new password are same")}

	ErrPermissionDenied = CustomError{
		Code: http.StatusForbidden,
		Err:  errors.New("permission denied"),
	}
)
