package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type TrWrapper struct {
	GormConn  *gorm.DB
	MongoConn mongo.SessionContext
}
