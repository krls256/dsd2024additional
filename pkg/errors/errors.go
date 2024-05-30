package errors

var (
	ErrEntityNotFound  = NewErrorWithCode("entity not found", 4000)
	ErrInternalError   = NewErrorWithCode("internal error", 0001)
	ErrDuplicateEntity = NewErrorWithCode("duplicate entity", 0002)
)
