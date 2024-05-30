package entities

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/errors"
)

type UUIDs []uuid.UUID

func (u *UUIDs) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return u.scanBytes(src)
	case string:
		return u.scanBytes([]byte(src))
	case nil:
		*u = nil

		return nil
	}

	return fmt.Errorf("%w: cannot convert %T to %T", errors.ErrInternalError, src, u)
}

func (u *UUIDs) scanBytes(src []byte) error {
	src = bytes.ReplaceAll(src, []byte{'{'}, []byte{})
	src = bytes.ReplaceAll(src, []byte{'}'}, []byte{})

	if len(src) == 0 {
		return nil
	}

	for _, item := range bytes.Split(src, []byte{','}) {
		id, err := uuid.Parse(string(item))
		if err != nil {
			return err
		}

		*u = append(*u, id)
	}

	return nil
}

func (u UUIDs) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}

	if n := len(u); n > 0 {
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = append(b, []byte(u[0].String())...)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = append(b, []byte(u[i].String())...)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}
