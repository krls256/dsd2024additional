package errs

import (
	"github.com/krls256/dsd2024additional/pkg/errors"
)

var (
	ErrAccountNotFound  = errors.NewErrorWithCode("account not found", 301000)
	ErrSessionsNotFound = errors.NewErrorWithCode("session not found", 304001)
	ErrTokenNotFound    = errors.NewErrorWithCode("token not found", 304002)
	ErrWrongPassword    = errors.NewErrorWithCode("wrong password", 301003)
)
