package dispatcher

import (
	"errors"
)

var (
	ErrUnknownAction = errors.New("unknown action")
	ErrInternal      = errors.New("conflict")
	ErrInvalidParams = errors.New("invalid parameters")
)
