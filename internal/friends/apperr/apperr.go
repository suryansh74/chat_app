package friends

import (
	"fmt"
)

var (
	ErrNotFound           = fmt.Errorf("friend not found")
	ErrAlreadyFriends     = fmt.Errorf("already friends")
	ErrUnauthorized       = fmt.Errorf("unauthorized")
	ErrInvalidRequestType = fmt.Errorf("invalid request type")
)
