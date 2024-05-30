package errs

import "github.com/krls256/dsd2024additional/pkg/errors"

var (
	ProfileAlreadyExists = errors.NewErrorWithCode("profile already exists", 401000)
	ProfileNotExists     = errors.NewErrorWithCode("profile not exists", 401001)
)
