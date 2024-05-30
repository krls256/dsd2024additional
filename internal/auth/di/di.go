package di

import (
	"github.com/krls256/dsd2024additional/internal/auth/constants"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/http"
	"github.com/krls256/dsd2024additional/internal/auth/migrations"
	"github.com/krls256/dsd2024additional/internal/auth/repositories"
	"github.com/krls256/dsd2024additional/internal/auth/repositories/pgsql"
	"github.com/krls256/dsd2024additional/internal/auth/repositories/redis"
	"github.com/krls256/dsd2024additional/internal/auth/rules"
	"github.com/krls256/dsd2024additional/internal/auth/services"
	"github.com/krls256/dsd2024additional/pkg/auth"
	pkgConstants "github.com/krls256/dsd2024additional/pkg/constants"
	pkgDI "github.com/krls256/dsd2024additional/pkg/di"
	"github.com/krls256/dsd2024additional/pkg/errors"
	pgsqlConn "github.com/krls256/dsd2024additional/pkg/pgsql"
	redisConn "github.com/krls256/dsd2024additional/pkg/redis"
	pkgRepositories "github.com/krls256/dsd2024additional/pkg/repositories/pgsql"
	"github.com/krls256/dsd2024additional/pkg/validator"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func PgSQLMigrations() pgsqlConn.SmartEmbed {
	return migrations.Migrations()
}

func Defs() []di.Def {
	return []di.Def{
		{
			Name: constants.AuthHandlerName,
			Tags: []di.Tag{{Name: pkgDI.HTTPHandlerTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				errorHandler := ctn.Get(pkgConstants.HTTPErrorHandlerName).(*errors.ErrorHTTPHandler)
				jwtFactory := ctn.Get(pkgConstants.JWTMiddlewareFactoryName).(*auth.JWTMiddlewareFactory)
				authService := ctn.Get(constants.AuthServiceName).(*services.AuthService)

				return http.NewAuthService(authService, errorHandler, jwtFactory), nil
			},
		},
		{
			Name: constants.AccountHandlerName,
			Tags: []di.Tag{{Name: pkgDI.HTTPHandlerTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				errorHandler := ctn.Get(pkgConstants.HTTPErrorHandlerName).(*errors.ErrorHTTPHandler)
				jwtFactory := ctn.Get(pkgConstants.JWTMiddlewareFactoryName).(*auth.JWTMiddlewareFactory)
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)

				return http.NewAccountHandler(accountService, errorHandler, jwtFactory), nil
			},
		},
		{
			Name: constants.AuthServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				authorizer := ctn.Get(pkgConstants.JWTAuthorizerName).(*auth.JWTAuthorizer)
				validatorEngine := ctn.Get(pkgConstants.ValidatorName).(*validator.Validator)

				tokenRepository := ctn.Get(constants.TokenRepositoryName).(*pkgRepositories.BaseRepository[*entities.Token])

				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)
				sessionService := ctn.Get(constants.SessionServiceName).(*services.SessionService)

				return services.NewAuthService(tokenRepository, accountService, sessionService,
					authorizer, validatorEngine), nil
			},
		},
		{
			Name: constants.SessionServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				sessionRepository := ctn.Get(constants.SessionRepositoryName).(repositories.SessionRepository)

				return services.NewSessionService(sessionRepository), nil
			},
		},
		{
			Name: constants.AccountServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				accountRepository := ctn.Get(constants.AccountRepositoryName).(*pgsql.AccountRepository)

				validatorEngine := ctn.Get(pkgConstants.ValidatorName).(*validator.Validator)

				return services.NewAccountService(accountRepository, validatorEngine), nil
			},
		},
		{
			Name: constants.AccountRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(pkgConstants.PgSQLName).(*gorm.DB)

				return pgsql.NewAccountRepository(conn), nil
			},
		},
		{
			Name: constants.TokenRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(pkgConstants.PgSQLName).(*gorm.DB)

				return pkgRepositories.NewBaseRepository[*entities.Token](conn), nil
			},
		},
		{
			Name: constants.SessionRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(pkgConstants.RedisName).(*redisConn.Client)

				return redis.NewSessionRepository(conn), nil
			},
		},
		{
			Name: constants.UniqueLoginCustomRuleName,
			Tags: []di.Tag{{Name: pkgDI.CustomRuleTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)

				return rules.NewUniqueLoginRule(accountService), nil
			},
		},
	}
}
