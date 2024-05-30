package migrations

import (
	"embed"
	pgsqlConn "github.com/krls256/dsd2024additional/pkg/pgsql"
)

//go:embed pgsql/*.sql
var embedMigrations embed.FS

func Migrations() pgsqlConn.SmartEmbed {
	return pgsqlConn.SmartEmbed{
		Embed:     embedMigrations,
		SourceDir: "pgsql",
	}
}
