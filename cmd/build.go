package cmd

import (
	authDI "github.com/krls256/dsd2024additional/internal/auth/di"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	"github.com/sarulabs/di/v2"
)

func Build() ([]di.Def, []pgsql.SmartEmbed) {
	pgSQLMigrations := []pgsql.SmartEmbed{
		authDI.PgSQLMigrations(),
	}

	defs := []di.Def{}
	defs = append(defs, authDI.Defs()...)

	return defs, pgSQLMigrations
}
