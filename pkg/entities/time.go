package entities

import (
	"database/sql/driver"
	"fmt"
	"github.com/krls256/dsd2024additional/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strings"
	"time"
)

type NanoTime struct {
	time.Time
}

func (n *NanoTime) Scan(src any) error {
	var err error

	switch src := src.(type) {
	case []byte:
		n.Time, err = time.Parse(string(src), time.RFC3339Nano)

		return err
	case string:
		n.Time, err = time.Parse(src, time.RFC3339Nano)

		return err
	case time.Time:
		n.Time = src

		return nil
	case *time.Time:
		n.Time = *src

		return nil
	case nil:
		return nil
	}

	return fmt.Errorf("%w: cannot convert %T to %T", errors.ErrInternalError, src, n)
}

func (n NanoTime) Value() (driver.Value, error) {
	return n.Time.Format(time.RFC3339Nano), nil
}

func (n *NanoTime) MarshalJSON() ([]byte, error) {
	if n.Time.IsZero() {
		return nil, nil
	}

	return []byte(fmt.Sprintf(`"%s"`, n.Time.Format(time.RFC3339Nano))), nil
}

func (n *NanoTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	if s == "null" {
		return nil
	}

	n.Time, err = time.Parse(time.RFC3339Nano, s)

	return
}

func (n *NanoTime) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	if err := bson.UnmarshalValue(btype, data, &n.Time); err != nil {
		return err
	}

	return nil
}

func (n *NanoTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(n.Time)
}
