package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/entities"
)

type BaseRepository[T entities.IdentifiableEntity] interface {
	Find(ctx context.Context, conditions map[string]interface{}) (data []T, err error)
	FindLimit(ctx context.Context, conditions map[string]interface{}, limit, offset int) (data []T, total int64, err error)
	FindBy(ctx context.Context, conditions map[string]interface{}, order []entities.Order) (T, bool, error)

	CreateNoReturn(ctx context.Context, m T) error
	SaveNoReturn(ctx context.Context, m T) error
	SaveMany(ctx context.Context, m []T) error

	Delete(ctx context.Context, ids []uuid.UUID, conditions map[string]interface{}) error
	DeleteAll(ctx context.Context, conditions map[string]interface{}) error

	CreateClosure(ctx context.Context, m T) func(tx TrWrapper) error
	SaveClosure(ctx context.Context, m T) func(tx TrWrapper) error
	SaveManyClosure(ctx context.Context, m []T) func(tx TrWrapper) error
	DeleteClosure(ctx context.Context, m T, conditions map[string]interface{}) func(tx TrWrapper) error
	DeleteAllClosure(ctx context.Context, conditions map[string]interface{}) func(tx TrWrapper) error

	Distinct(ctx context.Context, col string) ([]string, error)

	Paginate(ctx context.Context, filters map[string]interface{}, order []entities.Order, limit int, page int) (
		pagination entities.Pagination[T], err error)

	Count(ctx context.Context, conditions map[string]interface{}) (int64, error)
}
