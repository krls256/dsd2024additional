package mongo

import (
	"context"
	"errors"
	"github.com/krls256/dsd2024additional/pkg/entities"
	pkgErrors "github.com/krls256/dsd2024additional/pkg/errors"
	pkgMongo "github.com/krls256/dsd2024additional/pkg/mongo"
	"github.com/krls256/dsd2024additional/pkg/repositories"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository[T entities.IdentifiableEntity] struct {
	conn           *mongo.Database
	collectionName string
}

func NewBaseRepository[T entities.IdentifiableEntity](conn *mongo.Database, collectionName string) *BaseRepository[T] {
	return &BaseRepository[T]{
		conn:           conn,
		collectionName: collectionName,
	}
}

func (r *BaseRepository[T]) DB() *mongo.Database {
	return r.conn
}

func (r *BaseRepository[T]) Start() *mongo.Collection {
	return r.conn.Collection(r.collectionName)
}

func (r *BaseRepository[T]) Find(ctx context.Context, conditions map[string]interface{}) (data []T, err error) {
	cur, err := r.Start().Find(ctx, pkgMongo.MapToFilters(conditions))
	if err != nil {
		return nil, err
	}

	return pkgMongo.DecodeMany[T](ctx, cur)
}

func (r *BaseRepository[T]) FindLimit(ctx context.Context, conditions map[string]interface{}, limit, offset int) (
	data []T, total int64, err error) {
	filters := pkgMongo.MapToFilters(conditions)

	total, err = r.Start().CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cur, err := r.Start().Find(ctx, filters, findOptions)
	if err != nil {
		return nil, 0, err
	}

	data, err = pkgMongo.DecodeMany[T](ctx, cur)

	return data, total, err
}

func (r *BaseRepository[T]) FindBy(ctx context.Context, params map[string]interface{},
	order []entities.Order) (entity T, ok bool, err error) {
	findOptions := options.FindOne()
	findOptions.SetSort(pkgMongo.ToOrderMany(order))

	err = r.Start().FindOne(ctx, pkgMongo.MapToFilters(params), findOptions).Decode(&entity)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return entity, false, nil
	}

	if err != nil {
		return entity, false, err
	}

	return entity, true, nil
}

func (r *BaseRepository[T]) CreateNoReturn(ctx context.Context, entity T) error {
	_, err := r.Start().InsertOne(ctx, entity)

	if err != nil && mongo.IsDuplicateKeyError(err) {
		return pkgErrors.ErrDuplicateEntity
	}

	return err
}

func (r *BaseRepository[T]) Delete(ctx context.Context, ids []uuid.UUID, conditions map[string]interface{}) error {
	conditions["_id"] = ids

	_, err := r.Start().DeleteMany(ctx, pkgMongo.MapToFilters(conditions))

	return err
}

func (r *BaseRepository[T]) DeleteAll(ctx context.Context, conditions map[string]interface{}) error {
	_, err := r.Start().DeleteMany(ctx, pkgMongo.MapToFilters(conditions))

	return err
}

func (r *BaseRepository[T]) SaveNoReturn(ctx context.Context, entity T) error {
	opts := options.Update().SetUpsert(true)

	_, err := r.Start().UpdateOne(ctx, bson.M{
		"_id": entity.GetID(),
	}, bson.M{"$set": entity}, opts)

	return err
}

func (r *BaseRepository[T]) SaveMany(ctx context.Context, entities []T) error {
	upserts := []mongo.WriteModel{}

	for _, e := range entities {
		upserts = append(upserts, mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": e.GetID()}).
			SetUpdate(bson.M{"$set": e}).
			SetUpsert(true))
	}

	_, err := r.Start().BulkWrite(ctx, upserts)

	return err
}

func (r *BaseRepository[T]) CreateClosure(ctx context.Context, entity T) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		_, err := r.Start().InsertOne(tx.MongoConn, entity)

		return err
	}
}

func (r *BaseRepository[T]) SaveClosure(ctx context.Context, entity T) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		opts := options.Update().SetUpsert(true)

		_, err := r.Start().UpdateOne(tx.MongoConn, bson.M{
			"_id": entity.GetID(),
		}, bson.M{"$set": entity}, opts)

		return err
	}
}

func (r *BaseRepository[T]) SaveManyClosure(ctx context.Context, entities []T) func(tx repositories.TrWrapper) error {
	// Note: check correctness
	return func(tx repositories.TrWrapper) error {
		upserts := []mongo.WriteModel{}

		for _, e := range entities {
			upserts = append(upserts, mongo.NewUpdateManyModel().
				SetFilter(bson.M{"_id": e.GetID()}).
				SetUpdate(bson.M{"$set": e}).
				SetUpsert(true))
		}

		_, err := r.Start().BulkWrite(tx.MongoConn, upserts)

		return err
	}
}

func (r *BaseRepository[T]) DeleteClosure(ctx context.Context, entity T,
	conditions map[string]interface{}) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		conditions["_id"] = entity.GetID()

		_, err := r.Start().DeleteOne(tx.MongoConn, conditions)

		return err
	}
}

func (r *BaseRepository[T]) DeleteAllClosure(ctx context.Context,
	conditions map[string]interface{}) func(tx repositories.TrWrapper) error {
	return func(tx repositories.TrWrapper) error {
		_, err := r.Start().DeleteMany(ctx, pkgMongo.MapToFilters(conditions))

		return err
	}
}

func (r *BaseRepository[T]) Distinct(ctx context.Context, col string) ([]string, error) {
	resInterface, err := r.Start().Distinct(ctx, col, bson.M{})
	if err != nil {
		return nil, err
	}

	res := []string{}

	for _, item := range resInterface {
		str, ok := item.(string)
		if !ok {
			return nil, pkgErrors.ErrInternalError
		}

		res = append(res, str)
	}

	return res, nil
}

func (r *BaseRepository[T]) Paginate(ctx context.Context, conditions map[string]interface{}, order []entities.Order, limit int, page int) (
	pagination entities.Pagination[T], err error) {
	filters := pkgMongo.MapToFilters(conditions)

	total, err := r.Start().CountDocuments(ctx, filters)
	if err != nil {
		return pagination, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(limit * (page - 1)))
	findOptions.SetSort(pkgMongo.ToOrderMany(order))

	cur, err := r.Start().Find(ctx, filters, findOptions)
	if err != nil {
		return pagination, err
	}

	items, err := pkgMongo.DecodeMany[T](ctx, cur)
	if err != nil {
		return pagination, err
	}

	pagination.Total = int(total)
	pagination.Limit = limit
	pagination.CurrentPage = page
	pagination.Items = items

	return pagination, nil
}

func (r *BaseRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	filters := pkgMongo.MapToFilters(conditions)

	total, err := r.Start().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}

	return total, nil
}
