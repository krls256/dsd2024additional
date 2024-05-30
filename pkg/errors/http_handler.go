package errors

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
	"github.com/samber/lo"
)

var ErrInternal = errors.New("internal error")
var httpSemanticMap = map[int]func(ctx *fiber.Ctx, meta interface{}, err ...error) error{
	0:  http.BadRequest,
	1:  http.Unauthorized,
	2:  http.PaymentRequired,
	3:  http.Forbidden,
	4:  http.NotFound,
	9:  http.Conflict,
	22: http.ValidationFailed,
	60: http.ServerError,
}

type ErrorHTTPHandler struct {
	moduleCode    int
	internalError ErrorWithCode
}

func NewErrorHTTPHandler(moduleCode int) *ErrorHTTPHandler {
	return &ErrorHTTPHandler{
		moduleCode:    moduleCode,
		internalError: WrapErrorWithCode(ErrInternal, moduleCode),
	}
}

func (h *ErrorHTTPHandler) HandleError(ctx *fiber.Ctx, err error) error {
	codes := ErrorToCodes(err)
	if len(codes) == 0 {
		codes = append(codes, UnknownErrorCode(h.moduleCode))
	}

	semantics := lo.Map(codes, func(item int, index int) int {
		_, semantic, _ := SplitErrorCode(item)

		return semantic
	})

	fn, ok := httpSemanticMap[semantics[0]]
	if !ok {
		fn = http.ServerError
	}

	return fn(ctx, map[string]interface{}{
		"codes": codes,
	}, err)
}

func (h *ErrorHTTPHandler) HandleCustomError(ctx *fiber.Ctx, cr CustomErrorResponse, err error) error {
	if cr != nil {
		toResp, wrap := cr.Error(err)

		if wrap {
			return h.HandleError(ctx, fmt.Errorf("%w: %v", ErrInternal, err))
		}

		return http.OKRaw(ctx, toResp)
	}

	return h.HandleError(ctx, err)
}

type CustomErrorResponse interface {
	Error(err error) (data any, wrapInStandard bool)
}

func (h *ErrorHTTPHandler) HandleTaggedError(ctx *fiber.Ctx, err ErrorWithCode) error {
	code := err.Code()

	_, semantic, _ := SplitErrorCode(code)

	fn, ok := httpSemanticMap[semantic]
	if !ok {
		fn = http.ServerError
	}

	return fn(ctx, map[string]interface{}{
		"codes": []int{code},
	}, err)
}

func SplitErrorCode(code int) (module, semantic, exact int) {
	module = code / 100 / 1000
	semantic = code / 1000 % 100
	exact = code % 1000

	return
}

func JoinErrorCode(module, semantic, exact int) (code int) {
	return module*100*1000 + semantic*1000 + exact
}

func UnknownErrorCode(module int) (code int) {
	code = module*100000 + 60000

	return
}
