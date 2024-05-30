package pgsql

import (
	"context"
	"errors"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	pkgEntities "github.com/krls256/dsd2024additional/pkg/entities"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	pkgRepositories "github.com/krls256/dsd2024additional/pkg/repositories/pgsql"
	"gorm.io/gorm"
)

func NewAccountRepository(conn *gorm.DB) *AccountRepository {
	return &AccountRepository{
		BaseRepository: *pkgRepositories.NewBaseRepository[*entities.Account](conn),
	}
}

type AccountRepository struct {
	pkgRepositories.BaseRepository[*entities.Account]
}

func (r *AccountRepository) Find(ctx context.Context, conditions map[string]interface{}) (data []*entities.Account, err error) {
	query := r.DB().WithContext(ctx).Where(conditions)

	err = query.Find(&data).Error

	return data, err
}

func (r *AccountRepository) FindLimit(ctx context.Context, conditions map[string]interface{},
	limit, offset int) (data []*entities.Account, total int64, err error) {
	query := r.DB().WithContext(ctx).Where(conditions)

	query.Count(&total)

	err = query.Limit(limit).Offset(offset).Find(&data).Error

	return data, total, err
}
func (r *AccountRepository) FindBy(ctx context.Context, conditions map[string]interface{},
	order []pkgEntities.Order) (entity *entities.Account, ok bool, err error) {
	err = r.DB().WithContext(ctx).Where(conditions).Order(pgsql.ToOrderMany(order)).
		First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity, false, nil
	}

	if err != nil {
		return entity, false, err
	}

	return entity, true, nil
}
