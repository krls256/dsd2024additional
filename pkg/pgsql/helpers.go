package pgsql

import (
	"github.com/krls256/dsd2024additional/pkg/entities"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

func ToOrderOne(order entities.Order) clause.OrderByColumn {
	return clause.OrderByColumn{Column: clause.Column{Name: order.Column()}, Desc: order.IsDesc()}
}

func ToOrderMany(order []entities.Order) []clause.OrderByColumn {
	return lo.Map(order, func(item entities.Order, index int) clause.OrderByColumn {
		return ToOrderOne(item)
	})
}
