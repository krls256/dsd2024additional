package config

import (
	"errors"
	"fmt"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/mongo"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	"github.com/krls256/dsd2024additional/pkg/redis"
	"github.com/krls256/dsd2024additional/pkg/transport/http"

	"github.com/spf13/viper"
	"path/filepath"
)

var (
	ErrConfigParser = errors.New("config parser error")
)

type CustomConfig interface {
	Name() string
	ValuePtr() interface{}
}

type Config struct {
	HTTPConfig  http.Config
	PgSQLConfig *pgsql.Config
	MongoConfig *mongo.Config
	RedisConfig *redis.Config
	JWTConfig   *auth.JWTConfig

	CustomConfigs map[string]interface{}
}

func New(path string, containersConfig ...CustomConfig) (*Config, error) {
	viper.Reset()

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(filepath.Dir(abs))
	viper.SetConfigFile(filepath.Base(abs))

	config := &Config{
		CustomConfigs: map[string]interface{}{},
	}

	staticConfigs := map[string]interface{}{
		"http":  &config.HTTPConfig,
		"pgsql": &config.PgSQLConfig,
		"mongo": &config.MongoConfig,
		"redis": &config.RedisConfig,
		"jwt":   &config.JWTConfig,
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	for key, conf := range staticConfigs {
		if err := parseTagConfig(key, conf); err != nil {
			return nil, err
		}
	}

	for _, custom := range containersConfig {
		if _, ok := config.CustomConfigs[custom.Name()]; ok {
			return nil, fmt.Errorf("%w: custom key collision: %v", ErrConfigParser, custom.Name())
		}

		valuePtr := custom.ValuePtr()
		viperSub := viper.Sub(custom.Name())

		if err := parseSubConfig(viperSub, &valuePtr, custom.Name()); err != nil {
			return nil, err
		}

		config.CustomConfigs[custom.Name()] = valuePtr
	}

	return config, nil
}

func parseTagConfig(tag string, parseTo interface{}) error {
	subConfig := viper.Sub(tag)

	if err := parseSubConfig(subConfig, &parseTo, tag); err != nil {
		return err
	}

	return nil
}

func parseSubConfig(subConfig *viper.Viper, parseTo interface{}, name string) error {
	if subConfig == nil {
		return fmt.Errorf("%w: can not read %v config to %T: subconfig is nil", ErrConfigParser, name, parseTo)
	}

	if err := subConfig.Unmarshal(parseTo); err != nil {
		return err
	}

	return nil
}
