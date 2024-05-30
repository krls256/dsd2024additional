package mongo

import (
	"context"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/entities"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type Counter struct {
	TotalCount int `bson:"totalCount"`
}

type ID struct {
	ID uuid.UUID `bson:"_id"`
}

func ToMany[Entity any](items []Entity) []interface{} {
	return lo.Map(items, func(item Entity, index int) interface{} {
		return item
	})
}

func DecodeMany[Entity any](ctx context.Context, cursor *mongo.Cursor) ([]Entity, error) {
	res := []Entity{}

	for cursor.Next(ctx) {
		var elem Entity

		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}

		res = append(res, elem)
	}

	return res, nil
}

func MapToFilters(m map[string]interface{}) bson.M {
	res := bson.M{}

	if m == nil {
		return bson.M{}
	}

	for k, v := range m {
		if v == nil || IsNilPointer(v) {
			continue
		}

		defaultFilterHandler(res, k, v)
	}

	return res
}

func defaultFilterHandler(res map[string]interface{}, k string, v interface{}) {
	if isSlice, isEmpty := IsSlice(v); isSlice {
		if !isEmpty {
			res[k] = bson.M{"$in": v}
		}
	} else {
		res[k] = v
	}
}

func IsSlice(v interface{}) (isSlice, isEmpty bool) {
	isSlice = reflect.TypeOf(v).Kind() == reflect.Slice
	if isSlice {
		isEmpty = reflect.ValueOf(v).Len() == 0
	}

	return isSlice, isEmpty
}

func IsNilPointer(v interface{}) (isNil bool) {
	return reflect.TypeOf(v).Kind() == reflect.Pointer && reflect.ValueOf(v).IsNil()
}

func ToOrderMany(order []entities.Order) map[string]interface{} {
	res := map[string]interface{}{}

	for _, o := range order {
		val := 1

		if o.IsDesc() {
			val = -1
		}

		res[o.Column()] = val
	}

	return res
}
