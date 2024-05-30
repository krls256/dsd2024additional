package transport

import (
	"errors"
	"github.com/samber/lo"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Meta    interface{} `json:"meta"`
	Nonce   interface{} `json:"nonce,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(status int, meta, nonce, data interface{}) *Response {
	success := false
	if status >= 200 && status <= 299 {
		success = true
	}

	response := &Response{
		Status:  status,
		Success: success,
		Meta:    meta,
		Nonce:   nonce,
	}

	response.Data = StatusTextIfEmpty(data, status)

	return response
}

func NewFormattedError(status int, meta interface{}, nonce interface{}, errs []error) *Response {
	errs = lo.Filter(errs, func(item error, index int) bool {
		return item != nil
	})

	var err error
	if len(errs) == 0 {
		err = errors.New("no error")
	} else {
		err = errors.Join(errs...)
	}

	return NewResponse(status, meta, nonce, err.Error())
}

func StatusTextIfEmpty(data interface{}, status int) interface{} {
	if data == nil {
		return http.StatusText(status)
	}

	return data
}
