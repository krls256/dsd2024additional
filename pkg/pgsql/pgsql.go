package pgsql

import (
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SmartEmbed struct {
	Embed     embed.FS
	SourceDir string
}

func NewPgsqlConnection(config *Config, migrationEmbs ...SmartEmbed) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Pass, config.Name)

	log := logger.Default.LogMode(logger.Warn)
	if config.Silent {
		log = log.LogMode(logger.Silent)
	}

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log,
	})

	if err != nil {
		return nil, err
	}

	database, err := conn.DB()
	if err != nil {
		return nil, err
	}

	database.SetConnMaxIdleTime(config.MaxIdleTime)
	database.SetConnMaxLifetime(config.MaxLifetime)
	database.SetMaxIdleConns(config.MaxIdleConnections)
	database.SetMaxOpenConns(config.MaxConnections)

	if _, err = goose.EnsureDBVersion(database); err != nil {
		return nil, err
	}

	for _, migration := range migrationEmbs {
		goose.SetBaseFS(migration.Embed)

		if err := goose.Up(database, migration.SourceDir, goose.WithAllowMissing()); err != nil {
			return nil, err
		}
	}

	return conn, nil
}
