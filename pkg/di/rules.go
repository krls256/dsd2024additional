package di

import (
	"github.com/krls256/dsd2024additional/pkg/constants"
	"github.com/krls256/dsd2024additional/pkg/validator"
	"github.com/sarulabs/di/v2"
)

func BuildRules() []di.Def {
	return []di.Def{
		{
			Name: constants.NonZeroCustomRuleName,
			Tags: []di.Tag{{Name: CustomRuleTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.NewNonZeroTimeRule(), nil
			},
		},
		{
			Name: constants.BeforeNowCustomRuleName,
			Tags: []di.Tag{{Name: CustomRuleTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.NewBeforeNowRule(), nil
			},
		},
		{
			Name: constants.AfterNowCustomRuleName,
			Tags: []di.Tag{{Name: CustomRuleTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.NewAfterNowRule(), nil
			},
		},
	}
}
