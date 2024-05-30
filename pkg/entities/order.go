package entities

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrWrongOrderLen  = errors.New("order length must be equal to 2")
	ErrWrongOrderType = errors.New("order type must be on of: [desc, asc]")
)

var (
	OrderDecs = "desc"
	OrderAsc  = "asc"
)

type Order [2]string

func (o Order) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v %v"`, o[0], o[1])), nil
}

func (o *Order) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, " ")
	data = bytes.Trim(data, `"`)
	data = bytes.ReplaceAll(data, []byte{' ', ' '}, []byte{' '})

	split := bytes.Split(data, []byte{' '})

	const orderLenExpected = 2
	if len(split) != orderLenExpected {
		return fmt.Errorf("%w: but %v received", ErrWrongOrderLen, len(split))
	}

	orderType := string(split[1])
	if orderType != OrderDecs && orderType != OrderAsc {
		return fmt.Errorf("%w: but %v received", ErrWrongOrderType, orderType)
	}

	o[0], o[1] = string(split[0]), string(split[1])

	return nil
}

func (o Order) SetDesc() Order {
	o[1] = OrderDecs

	return o
}

func (o Order) SetAsc() Order {
	o[1] = OrderAsc

	return o
}

func (o Order) SetColumn(col string) Order {
	o[0] = col

	return o
}

func (o Order) IsDesc() bool {
	return o[1] == OrderDecs
}

func (o Order) IsAsc() bool {
	return o[1] == OrderAsc
}

func (o Order) Column() string {
	return o[0]
}

func (o Order) ToSlice() []Order {
	return []Order{o}
}
