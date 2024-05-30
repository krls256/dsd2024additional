package pgsql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/entities"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	"github.com/krls256/dsd2024additional/pkg/repositories"
	"gorm.io/gorm"
)

type BaseRepository[T entities.IdentifiableEntity] struct {
	conn *gorm.DB
}

func NewBaseRepository[T entities.IdentifiableEntity](conn *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		conn: conn,
	}
}

func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.conn
}

func (r *BaseRepository[T]) Find(ctx context.Context, conditions map[string]interface{}) (data []T, err error) {
	query := r.conn.WithContext(ctx).Where(conditions)

	err = query.Find(&data).Error

	return data, err
}

func (r *BaseRepository[T]) FindLimit(ctx context.Context, conditions map[string]interface{},
	limit, offset int) (data []T, total int64, err error) {
	query := r.conn.WithContext(ctx).Where(conditions)

	query.Count(&total)

	err = query.Limit(limit).Offset(offset).Find(&data).Error

	return data, total, err
}

func (r *BaseRepository[T]) FindBy(ctx context.Context, conditions map[string]interface{},
	order []entities.Order) (entity T, ok bool, err error) {
	err = r.conn.WithContext(ctx).Where(conditions).Order(pgsql.ToOrderMany(order)).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity, false, nil
	}

	if err != nil {
		return entity, false, err
	}

	return entity, true, nil
}

func (r *BaseRepository[T]) CreateNoReturn(ctx context.Context, m T) error {
	return r.conn.WithContext(ctx).Create(&m).Error
}
func (r *BaseRepository[T]) SaveNoReturn(ctx context.Context, m T) error {
	return r.conn.WithContext(ctx).Save(&m).Error
}

func (r *BaseRepository[T]) SaveMany(ctx context.Context, m []T) error {
	return r.conn.WithContext(ctx).Save(m).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, ids []uuid.UUID, conditions map[string]interface{}) error {
	var entity T

	return r.conn.WithContext(ctx).Model(&entity).Where("id in (?)", ids).Delete(nil, conditions).Error
}

func (r *BaseRepository[T]) DeleteAll(ctx context.Context, conditions map[string]interface{}) error {
	var entity T

	return r.conn.WithContext(ctx).Model(&entity).Delete(nil, conditions).Error
}

func (r *BaseRepository[T]) CreateClosure(ctx context.Context, m T) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		return tx.GormConn.WithContext(ctx).Create(&m).Error
	}
}
func (r *BaseRepository[T]) SaveClosure(ctx context.Context, m T) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		return tx.GormConn.WithContext(ctx).Save(&m).Error
	}
}

func (r *BaseRepository[T]) SaveManyClosure(ctx context.Context, m []T) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		return tx.GormConn.WithContext(ctx).Save(m).Error
	}
}

func (r *BaseRepository[T]) DeleteClosure(ctx context.Context, m T,
	conditions map[string]interface{}) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		return tx.GormConn.WithContext(ctx).Model(&m).Delete(&m, conditions).Error
	}
}

func (r *BaseRepository[T]) DeleteAllClosure(ctx context.Context,
	conditions map[string]interface{}) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		var entity T

		return tx.GormConn.WithContext(ctx).Model(&entity).Delete(nil, conditions).Error
	}
}

func (r *BaseRepository[T]) Distinct(ctx context.Context, col string) ([]string, error) {
	var (
		e     T
		items []string
	)

	return items, r.conn.WithContext(ctx).Model(&e).Distinct(col).Find(&items).Error
}

func (r *BaseRepository[T]) Paginate(ctx context.Context, filters map[string]interface{}, order []entities.Order, limit int, page int) (
	pagination entities.Pagination[T], err error) {
	conn := r.conn.WithContext(ctx).Where(filters).Order(pgsql.ToOrderMany(order))

	items := make([]T, 0)

	var entity T

	var total int64

	if err = conn.Model(&entity).Count(&total).Error; err != nil {
		return pagination, err
	}

	conn = conn.
		Limit(limit).
		Offset(limit * (page - 1))

	if err = conn.Find(&items).Error; err != nil {
		return
	}

	pagination.Total = int(total)
	pagination.Limit = limit
	pagination.CurrentPage = page
	pagination.Items = items

	return
}

func (r *BaseRepository[T]) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	conn := r.conn.WithContext(ctx).Where(filters)

	var (
		entity T
		total  int64
	)

	if err := conn.Model(&entity).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
