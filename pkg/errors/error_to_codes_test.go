package errors

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

// go test ./pkg/errs -run Test_ErrorToCodes -v
func Test_ErrorToCodes(t *testing.T) {
	err1 := NewErrorWithCode("err1", 1)
	err2 := NewErrorWithCode("err2", 2)

	err3 := fmt.Errorf("%w: err3", err1)
	err4 := fmt.Errorf("%w: err4", err2)

	err5 := errors.Join(err3, err4)

	err6 := fmt.Errorf("%w: err6", err5)

	require.Equal(t, []int{1, 2}, ErrorToCodes(err6))
	require.Equal(t, []int{1, 2}, ErrorToCodes(err5))
	require.Equal(t, []int{2}, ErrorToCodes(err4))
	require.Equal(t, []int{1}, ErrorToCodes(err3))
	require.Equal(t, []int{2}, ErrorToCodes(err2))
	require.Equal(t, []int{1}, ErrorToCodes(err1))
}
