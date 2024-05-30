package di

import (
	"errors"
	"github.com/krls256/dsd2024additional/pkg/constants"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	"github.com/krls256/dsd2024additional/pkg/validator"

	"github.com/sarulabs/di/v2"
)

var ErrInvalidBuild = errors.New("invalid build")

func Build(configPath string, pgSQLMigrations []pgsql.SmartEmbed, defs ...di.Def) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	defs = append(defs, PkgDefs(configPath, pgSQLMigrations)...)
	defs = append(defs, BuildRules()...)

	if err := builder.Add(defs...); err != nil {
		panic(err)
	}

	ctn := builder.Build()

	v, ok := ctn.Get(constants.ValidatorName).(*validator.Validator)
	if !ok {
		return nil, ErrInvalidBuild
	}

	customRulesNames := filterNamesByTag(ctn.Definitions(), CustomRuleTag)

	customRules := []validator.CustomRule{}

	for _, name := range customRulesNames {
		customRules = append(customRules, ctn.Get(name).(validator.CustomRule))
	}

	if err := v.AddCustomRules(customRules...); err != nil {
		return nil, err
	}

	return ctn, nil
}
