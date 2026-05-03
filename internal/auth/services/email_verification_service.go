package auth

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/suryansh74/chat_app/internal/auth/repositories"
)

var (
	ErrInvalidOTP      = errors.New("invalid OTP")
	ErrOTPExpired      = errors.New("OTP has expired")
	ErrMaxAttempts     = errors.New("max attempts exceeded")
	ErrAlreadyVerified = errors.New("email already verified")
)

type EmailVerificationService struct {
	repo             repositories.UserRepository
	otpExpiryMinutes int
	otpMaxAttempts   int
}

func NewEmailVerificationService(repo repositories.UserRepository, otpExpiryMinutes, otpMaxAttempts int) *EmailVerificationService {
	return &EmailVerificationService{
		repo:             repo,
		otpExpiryMinutes: otpExpiryMinutes,
		otpMaxAttempts:   otpMaxAttempts,
	}
}

func (s *EmailVerificationService) GenerateOTP() (string, error) {
	const digits = "0123456789"
	result := make([]byte, 6)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = digits[n.Int64()]
	}
	return string(result), nil
}

func (s *EmailVerificationService) SendOTP(email string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user.IsVerified {
		return "", ErrAlreadyVerified
	}

	otp, err := s.GenerateOTP()
	if err != nil {
		return "", err
	}

	user.OTP = otp
	user.OTPExpiry = time.Now().Add(time.Duration(s.otpExpiryMinutes) * time.Minute)
	user.OTPAttempts = 0

	if err := s.repo.UpdateUser(user); err != nil {
		return "", err
	}

	return otp, nil
}

func (s *EmailVerificationService) VerifyOTP(email, otp string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return ErrAlreadyVerified
	}

	if user.OTP == "" {
		return ErrInvalidOTP
	}

	if time.Now().After(user.OTPExpiry) {
		return ErrOTPExpired
	}

	user.OTPAttempts++
	if user.OTPAttempts >= s.otpMaxAttempts {
		user.OTP = ""
		user.OTPExpiry = time.Time{}
		user.OTPAttempts = 0
		s.repo.UpdateUser(user)
		return ErrMaxAttempts
	}

	if user.OTP != otp {
		s.repo.UpdateUser(user)
		return ErrInvalidOTP
	}

	user.IsVerified = true
	user.OTP = ""
	user.OTPExpiry = time.Time{}
	user.OTPAttempts = 0

	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *EmailVerificationService) IsVerified(email string) (bool, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user.IsVerified, nil
}
