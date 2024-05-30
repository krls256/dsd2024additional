package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const InitWaitTimeout = 10 * time.Second

func NewMongoConnection(config *Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), InitWaitTimeout)
	defer cancel()

	// monitor := &event.CommandMonitor{
	//	Started: func(_ context.Context, e *event.CommandStartedEvent) {
	//		if e.CommandName != "endSessions" {
	//			zap.S().Info(e.Command.String())
	//		}
	//	},
	//}

	url := fmt.Sprintf("mongodb://%v:%v@%v:%v", config.User, config.Pass, config.Host, config.Port)
	if config.IsEmptyUser() {
		url = fmt.Sprintf("mongodb://%v:%v", config.Host, config.Port)
	}

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(url))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
