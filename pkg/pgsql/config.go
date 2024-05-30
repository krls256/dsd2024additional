package pgsql

import "time"

type Config struct {
	Host               string
	Port               uint16
	Name               string
	User               string
	Pass               string
	MaxIdleTime        time.Duration
	MaxLifetime        time.Duration
	Compression        string
	MaxIdleConnections int
	MaxConnections     int
	Silent             bool
}
