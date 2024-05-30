package di

import (
	"github.com/krls256/dsd2024additional/internal/web/constants"
	"github.com/krls256/dsd2024additional/internal/web/http"
	pkgDI "github.com/krls256/dsd2024additional/pkg/di"
	"github.com/sarulabs/di/v2"
)

func Defs() []di.Def {
	return []di.Def{
		{
			Name: constants.WebHandlerName,
			Tags: []di.Tag{{Name: pkgDI.HTTPHandlerTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				return http.NewWebHandler(), nil
			},
		},
	}
}
