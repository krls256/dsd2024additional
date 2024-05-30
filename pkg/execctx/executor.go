package execctx

import (
	"context"
	pkgMongo "github.com/krls256/dsd2024additional/pkg/mongo"
	"github.com/krls256/dsd2024additional/pkg/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

func NewExecutorFactory(mongoConn *mongo.Client, mongCfg *pkgMongo.Config, gormConn *gorm.DB) *ExecutorFactory {
	return &ExecutorFactory{mongoConn: mongoConn, mongCfg: mongCfg, gormConn: gormConn}
}

type ExecutorFactory struct {
	mongoConn *mongo.Client
	mongCfg   *pkgMongo.Config
	gormConn  *gorm.DB
}

func (f *ExecutorFactory) Executor() *Executor {
	return &Executor{
		mongoConn: f.mongoConn,
		mongCfg:   f.mongCfg,
		gormConn:  f.gormConn,
	}
}

type Executor struct {
	mongoConn *mongo.Client
	mongCfg   *pkgMongo.Config
	gormConn  *gorm.DB

	once sync.Once

	operationsPool []func(tx repositories.TrWrapper) error
}

func (r *Executor) Add(op func(tx repositories.TrWrapper) error) {
	r.operationsPool = append(r.operationsPool, op)
}

// Exec It's not a atomic operation for mongo db till replica set
func (r *Executor) Exec(ctx context.Context) (executed bool, err error) {
	r.once.Do(func() {
		executed = true

		mongoSession, localErr := r.mongoConn.StartSession()
		if localErr != nil {
			err = localErr

			return
		}

		defer mongoSession.EndSession(ctx)

		mongoTx := mongo.NewSessionContext(ctx, mongoSession)

		gormTx := r.gormConn.Begin()

		tx := repositories.TrWrapper{GormConn: gormTx, MongoConn: mongoTx}

		for _, op := range r.operationsPool {
			if txErr := op(tx); txErr != nil {
				rollbackGorm(gormTx)

				err = txErr

				return
			}
		}

		if txErr := gormTx.Commit().Error; txErr != nil {
			err = txErr

			return
		}
	})

	return executed, err
}

func rollbackGorm(gormTx *gorm.DB) {
	if err := gormTx.Rollback().Error; err != nil {
		zap.S().Error(err)
	}
}
