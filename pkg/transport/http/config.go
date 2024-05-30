package http

import (
	"fmt"
	"github.com/krls256/dsd2024additional/pkg/transport/http/middlewares"
	"time"
)

type Config struct {
	Host         string
	Port         uint16
	Silent       bool
	LoggerConfig middlewares.LoggerConfig
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (c Config) DNS() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
