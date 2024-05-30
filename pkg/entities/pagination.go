package entities

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	"time"
)

type Entity interface {
	IdentifiableEntity
	PreCreateCompute() error
	PreUpdateCompute(UpdateEntityRequest) error
}

type UpdateEntityRequest interface {
	IdentifiableEntity
}

type DeleteEntityRequest interface {
	IdentifiableEntity
	SetableIDEntity
}

type IdentifiableEntity interface {
	GetID() uuid.UUID
	IDColumnName() string
}

type SetableIDEntity interface {
	SetID(id uuid.UUID)
}

type ITimeFilter interface {
	From() (NanoTime, bool)
	To() (NanoTime, bool)
}

type PaginationFilters interface {
	Filters() (map[string]interface{}, error)
}

type CRUDPaginationRequest[Filters PaginationFilters] struct {
	Page    int     `json:"page" form:"page" validate:"gte=1"`
	Limit   int     `json:"limit" form:"limit" validate:"gte=1"`
	Order   []Order `json:"order"`
	Filters Filters `json:"filters"`
}

type PaginationRequest[Filters any] struct {
	Page    int     `json:"page" form:"page" validate:"gte=1"`
	Limit   int     `json:"limit" form:"limit" validate:"gte=1"`
	Order   []Order `json:"order"`
	Filters Filters `json:"filters"`
}

type Pagination[T any] struct {
	Items       []T `json:"items"`
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
}

type CurrencyGetter interface {
	GetCurrency() string
}

func EntityCurrencies[T CurrencyGetter](items []T) []string {
	currencies := lo.Map(items, func(item T, index int) string {
		return item.GetCurrency()
	})

	return lo.Uniq(currencies)
}

type TimestampGetter interface {
	GetTimestamp() time.Time
}

func StartEndTimestamps[T TimestampGetter](items []T) (start, end time.Time) {
	if len(items) == 0 {
		return start, end
	}

	start, end = items[0].GetTimestamp(), items[0].GetTimestamp()

	for i := 1; i < len(items); i++ {
		item := items[i].GetTimestamp()
		if item.Before(start) {
			start = item
		}

		if item.After(end) {
			end = item
		}
	}

	return start, end
}
