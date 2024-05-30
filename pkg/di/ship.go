package di

import (
	"context"
	"fmt"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/config"
	"github.com/krls256/dsd2024additional/pkg/constants"
	"github.com/krls256/dsd2024additional/pkg/errors"
	"github.com/krls256/dsd2024additional/pkg/execctx"
	pkgMongo "github.com/krls256/dsd2024additional/pkg/mongo"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	"github.com/krls256/dsd2024additional/pkg/redis"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
	"github.com/krls256/dsd2024additional/pkg/transport/http/middlewares"
	"github.com/krls256/dsd2024additional/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func PkgDefs(configPath string, migrations []pgsql.SmartEmbed) []di.Def {
	return []di.Def{
		{
			Name: constants.LoggerName,
			Build: func(ctn di.Container) (interface{}, error) {
				logger, err := zap.NewDevelopment()

				if err != nil {
					return nil, fmt.Errorf("can't initialize zap logger: %v", err)
				}

				zap.ReplaceGlobals(logger)

				return logger, nil
			},
		},
		{
			Name: constants.ConfigName,
			Build: func(ctn di.Container) (interface{}, error) {
				customConfigNames := filterNamesByTag(ctn.Definitions(), CustomConfigTag)

				customConfig := []config.CustomConfig{}

				for _, name := range customConfigNames {
					customConfig = append(customConfig, ctn.Get(name).(config.CustomConfig))
				}

				return config.New(configPath, customConfig...)
			},
		},
		{
			Name: constants.ValidatorName,
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.New()
			},
		},
		{
			Name: constants.ExecutorFactoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				mongoConn := ctn.Get(constants.MongoClientName).(*mongo.Client)
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				gormConn := ctn.Get(constants.PgSQLName).(*gorm.DB)

				return execctx.NewExecutorFactory(mongoConn, cfg.MongoConfig, gormConn), nil
			},
		},
		{
			Name: constants.ExecutorCtxFactoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				execFactory := ctn.Get(constants.ExecutorFactoryName).(*execctx.ExecutorFactory)

				return execctx.NewContextFactory(execFactory), nil
			},
		},
		{
			Name: constants.PgSQLName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return pgsql.NewPgsqlConnection(cfg.PgSQLConfig, migrations...)
			},
		},
		{
			Name: constants.MongoClientName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return pkgMongo.NewMongoConnection(cfg.MongoConfig)
			},
		},
		{
			Name: constants.MongoDatabaseName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				mongoConn := ctn.Get(constants.MongoClientName).(*mongo.Client)

				return mongoConn.Database(cfg.MongoConfig.Name), nil
			},
		},
		{
			Name: constants.RedisName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return redis.New(cfg.RedisConfig)
			},
		},
		{
			Name: constants.JWTAuthorizerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return auth.NewAuthorizer(cfg.JWTConfig), nil
			},
		},
		{
			Name: constants.LogMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return middlewares.RequestLoggerMiddleware(cfg.HTTPConfig.LoggerConfig), nil
			},
		},
		{
			Name: constants.HTTPServerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				handlerNames := filterNamesByTag(ctn.Definitions(), HTTPHandlerTag)

				handlers := []http.Handler{}

				for _, name := range handlerNames {
					handlers = append(handlers, ctn.Get(name).(http.Handler))
				}

				mws := []fiber.Handler{
					ctn.Get(constants.LogMiddlewareName).(fiber.Handler),
				}

				return http.NewServer(context.Background(), "chat", cfg.HTTPConfig, handlers, mws), nil
			},
		},
		{
			Name: constants.SessionServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.RedisName).(*redis.Client)

				return auth.NewSessionService(conn), nil
			},
		},
		{
			Name: constants.HTTPErrorHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				return errors.NewErrorHTTPHandler(1), nil
			},
		},
		{
			Name: constants.JWTMiddlewareFactoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				authorizer := ctn.Get(constants.JWTAuthorizerName).(*auth.JWTAuthorizer)
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				ss := ctn.Get(constants.SessionServiceName).(*auth.SessionService)

				return auth.NewJWTMiddlewareFactory(authorizer, cfg.JWTConfig, ss), nil
			},
		},
	}
}
