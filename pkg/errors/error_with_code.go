package errors

import (
	"errors"
	"github.com/samber/lo"
)

type ErrorWithCode struct {
	err  error
	code int
}

func NewErrorWithCode(err string, code int) ErrorWithCode {
	return ErrorWithCode{
		err:  errors.New(err),
		code: code,
	}
}

func WrapErrorWithCode(err error, code int) ErrorWithCode {
	return ErrorWithCode{
		err:  err,
		code: code,
	}
}

func (e ErrorWithCode) Error() string {
	return e.err.Error()
}

func (e ErrorWithCode) Code() int {
	return e.code
}

type manyUnwrap interface {
	Unwrap() []error
}

func ErrorToCodes(err error) (codes []int) {
	errorToCodes(err, &codes)

	return lo.Uniq(codes)
}

func errorToCodes(err error, codes *[]int) {
	for err != nil {
		var ewc ErrorWithCode
		if ok := errors.As(err, &ewc); ok {
			*codes = append(*codes, ewc.Code())
		}

		var mu manyUnwrap
		if ok := errors.As(err, &mu); ok {
			for _, e := range mu.Unwrap() {
				errorToCodes(e, codes)
			}

			return
		}

		err = errors.Unwrap(err)
	}
}

type ErrorWithCodes interface {
	error
	Codes() []int
}

func WrapErrorWithCodes(err error, codes []int) ErrorWithCodes {
	return errorWithCodes{
		err:   err,
		codes: codes,
	}
}

type errorWithCodes struct {
	err   error
	codes []int
}

func (e errorWithCodes) Error() string {
	return e.err.Error()
}

func (e errorWithCodes) Codes() []int {
	return e.codes
}
