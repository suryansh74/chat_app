package notification

import (
	"fmt"
)

var (
	ErrNotFound           = fmt.Errorf("notification not found")
	ErrUnauthorized       = fmt.Errorf("unauthorized")
	ErrInvalidRequestType = fmt.Errorf("invalid request type")
)
