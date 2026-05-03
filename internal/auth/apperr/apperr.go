package apperr

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidResetToken  = errors.New("invalid reset token")
	ErrResetTokenExpired  = errors.New("reset token has expired")
	ErrMaxResetAttempts   = errors.New("max reset attempts exceeded")
	ErrPasswordTooWeak    = errors.New("password is too weak")
	ErrPasswordMismatch   = errors.New("password and confirmation do not match")
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

func IsInvalidResetToken(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == "INVALID_RESET_TOKEN"
	}
	return errors.Is(err, ErrInvalidResetToken)
}

func NewInvalidResetToken() *AppError {
	return &AppError{
		Code:    "INVALID_RESET_TOKEN",
		Message: "invalid reset token",
		Err:     ErrInvalidResetToken,
	}
}

func IsResetTokenExpired(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == "RESET_TOKEN_EXPIRED"
	}
	return errors.Is(err, ErrResetTokenExpired)
}

func NewResetTokenExpired() *AppError {
	return &AppError{
		Code:    "RESET_TOKEN_EXPIRED",
		Message: "reset token has expired",
		Err:     ErrResetTokenExpired,
	}
}

func IsMaxResetAttempts(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == "MAX_RESET_ATTEMPTS_EXCEEDED"
	}
	return errors.Is(err, ErrMaxResetAttempts)
}

func NewMaxResetAttempts() *AppError {
	return &AppError{
		Code:    "MAX_RESET_ATTEMPTS_EXCEEDED",
		Message: "max reset attempts exceeded",
		Err:     ErrMaxResetAttempts,
	}
}

func IsPasswordTooWeak(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == "PASSWORD_TOO_WEAK"
	}
	return errors.Is(err, ErrPasswordTooWeak)
}

func NewPasswordTooWeak() *AppError {
	return &AppError{
		Code:    "PASSWORD_TOO_WEAK",
		Message: "password is too weak",
		Err:     ErrPasswordTooWeak,
	}
}

func IsPasswordMismatch(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == "PASSWORD_MISMATCH"
	}
	return errors.Is(err, ErrPasswordMismatch)
}

func NewPasswordMismatch() *AppError {
	return &AppError{
		Code:    "PASSWORD_MISMATCH",
		Message: "password and confirmation do not match",
		Err:     ErrPasswordMismatch,
	}
}
