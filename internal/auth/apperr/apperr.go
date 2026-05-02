package apperr

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func IsEmailAlreadyExists(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == "EMAIL_ALREADY_EXISTS"
	}
	return errors.Is(err, ErrEmailAlreadyExists)
}

func NewEmailAlreadyExists(email string) *AppError {
	return &AppError{
		Code:    "EMAIL_ALREADY_EXISTS",
		Message: "email " + email + " already exists",
		Err:     ErrEmailAlreadyExists,
	}
}
