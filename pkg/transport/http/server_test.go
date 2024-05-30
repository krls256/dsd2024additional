package http

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	correctCfg = Config{
		Host:         "0.0.0.0",
		Port:         80,
		Silent:       true,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}
)

// go test ./pkg/transport/http -run Test_Server_DoubleRunShutdown
func Test_Server_DoubleRunShutdown(t *testing.T) {
	server := NewServer(context.Background(), "", correctCfg, nil, nil)

	server.AsyncRun()

	time.Sleep(time.Millisecond * 100)

	require.Error(t, server.Run())
	require.NoError(t, server.Shutdown(context.Background()))
}

// go test ./pkg/transport/http -run Test_Server_GracefulShutdown_CanceledCtx
func Test_Server_GracefulShutdown_OutdatedCtx(t *testing.T) {
	server := NewServer(context.Background(), "", correctCfg, nil, nil)

	server.AsyncRun()

	time.Sleep(time.Millisecond * 100)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	time.Sleep(time.Millisecond * 100)

	require.NoError(t, server.Shutdown(ctx))
}
