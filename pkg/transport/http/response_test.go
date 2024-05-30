package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func newOneRouteHandler(path string, handler fiber.Handler) *oneRouteHandler {
	return &oneRouteHandler{path: path, handler: handler}
}

type oneRouteHandler struct {
	path    string
	handler fiber.Handler
}

func (h *oneRouteHandler) Register(router fiber.Router) {
	router.Get(h.path, h.handler)
}

func (h *oneRouteHandler) GenerateRequest(t *testing.T, cfg Config) *http.Request {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+filepath.Join(cfg.DNS(), h.path), nil)
	require.Nil(t, err)

	return r
}

// go test ./pkg/transport/http -run Test_Response_OK
func Test_Response_OK(t *testing.T) {
	h := newOneRouteHandler("test", func(ctx *fiber.Ctx) error {
		return OK(ctx, nil, nil)
	})

	server := NewServer(context.Background(), "", correctCfg, []Handler{h}, nil)
	server.AsyncRun()

	time.Sleep(time.Millisecond * 100)

	defer func() {
		require.Nil(t, server.Shutdown(context.Background()))
	}()

	req := h.GenerateRequest(t, correctCfg)
	resp, err := http.DefaultClient.Do(req)

	resp.Body.Close()

	require.Nil(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
}

// go test ./pkg/transport/http -run Test_Response_BadRequest
func Test_Response_BadRequest(t *testing.T) {
	testResponseWriter(t, BadRequest, http.StatusBadRequest)
}

// go test ./pkg/transport/http -run Test_Response_Unauthorized
func Test_Response_Unauthorized(t *testing.T) {
	testResponseWriter(t, Unauthorized, http.StatusUnauthorized)
}

// go test ./pkg/transport/http -run Test_Response_PaymentRequired
func Test_Response_PaymentRequired(t *testing.T) {
	testResponseWriter(t, PaymentRequired, http.StatusPaymentRequired)
}

// go test ./pkg/transport/http -run Test_Response_Forbidden
func Test_Response_Forbidden(t *testing.T) {
	testResponseWriter(t, Forbidden, http.StatusForbidden)
}

// go test ./pkg/transport/http -run Test_Response_NotFound
func Test_Response_NotFound(t *testing.T) {
	testResponseWriter(t, NotFound, http.StatusNotFound)
}

// go test ./pkg/transport/http -run Test_Response_Conflict
func Test_Response_Conflict(t *testing.T) {
	testResponseWriter(t, Conflict, http.StatusConflict)
}

// go test ./pkg/transport/http -run Test_Response_Teapot
func Test_Response_Teapot(t *testing.T) {
	testResponseWriter(t, Teapot, http.StatusTeapot)
}

// go test ./pkg/transport/http -run Test_Response_ValidationFailed
func Test_Response_ValidationFailed(t *testing.T) {
	testResponseWriter(t, ValidationFailed, http.StatusUnprocessableEntity)
}

// go test ./pkg/transport/http -run Test_Response_TooEarly
func Test_Response_TooEarly(t *testing.T) {
	testResponseWriter(t, TooEarly, http.StatusTooEarly)
}

// go test ./pkg/transport/http -run Test_Response_TooManyRequests
func Test_Response_TooManyRequests(t *testing.T) {
	testResponseWriter(t, TooManyRequests, http.StatusTooManyRequests)
}

// go test ./pkg/transport/http -run Test_Response_ServerError
func Test_Response_ServerError(t *testing.T) {
	testResponseWriter(t, ServerError, http.StatusInternalServerError)
}

func testResponseWriter(t *testing.T, fn func(ctx *fiber.Ctx, meta interface{}, err ...error) error, status int) {
	h := newOneRouteHandler("test", func(ctx *fiber.Ctx) error {
		return fn(ctx, nil)
	})

	server := NewServer(context.Background(), "", correctCfg, []Handler{h}, nil)
	server.AsyncRun()

	time.Sleep(time.Millisecond * 100)

	defer func() {
		require.Nil(t, server.Shutdown(context.Background()))
	}()

	req := h.GenerateRequest(t, correctCfg)
	resp, err := http.DefaultClient.Do(req)

	require.Nil(t, err)
	require.Equal(t, resp.StatusCode, status)

	defer resp.Body.Close()
}
