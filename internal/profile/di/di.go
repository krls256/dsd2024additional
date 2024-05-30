package di

import (
	"github.com/krls256/dsd2024additional/internal/profile/constants"
	"github.com/krls256/dsd2024additional/internal/profile/entities"
	"github.com/krls256/dsd2024additional/internal/profile/http"
	"github.com/krls256/dsd2024additional/internal/profile/migrations"
	"github.com/krls256/dsd2024additional/internal/profile/services"
	"github.com/krls256/dsd2024additional/pkg/auth"
	pkgConstants "github.com/krls256/dsd2024additional/pkg/constants"
	pkgDI "github.com/krls256/dsd2024additional/pkg/di"
	"github.com/krls256/dsd2024additional/pkg/errors"
	pgsqlConn "github.com/krls256/dsd2024additional/pkg/pgsql"
	pkgRepositories "github.com/krls256/dsd2024additional/pkg/repositories/pgsql"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func PgSQLMigrations() pgsqlConn.SmartEmbed {
	return migrations.Migrations()
}

func Defs() []di.Def {
	return []di.Def{
		{
			Name: constants.ProfileRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(pkgConstants.PgSQLName).(*gorm.DB)

				return pkgRepositories.NewBaseRepository[*entities.Profile](conn), nil
			},
		},
		{
			Name: constants.ProfileServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				repo := ctn.Get(constants.ProfileRepositoryName).(*pkgRepositories.BaseRepository[*entities.Profile])

				return services.NewProfileService(repo), nil
			},
		},
		{
			Name: constants.ProfileHandlerName,
			Tags: []di.Tag{{Name: pkgDI.HTTPHandlerTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				errorHandler := ctn.Get(pkgConstants.HTTPErrorHandlerName).(*errors.ErrorHTTPHandler)
				jwtFactory := ctn.Get(pkgConstants.JWTMiddlewareFactoryName).(*auth.JWTMiddlewareFactory)
				profileService := ctn.Get(constants.ProfileServiceName).(*services.ProfileService)

				return http.NewProfileHandler(profileService, errorHandler, jwtFactory), nil
			},
		},
	}
}
