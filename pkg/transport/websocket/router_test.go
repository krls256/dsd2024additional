package websocket

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Router_Accept(t *testing.T) {
	r := NewRouter()

	r.Accept("test", func(ctx *Context) {})

	require.Equal(t, 1, len(r.routeMap))
}

func Test_Router_FindExisting(t *testing.T) {
	r := NewRouter()

	r.Accept("test", func(ctx *Context) {})

	fn, ok := r.find("test")

	require.Equal(t, true, ok)
	require.NotNil(t, fn)
}

func Test_Router_FindNotExisting(t *testing.T) {
	r := NewRouter()

	fn, ok := r.find("test")

	require.Equal(t, false, ok)
	require.Nil(t, fn)
}
