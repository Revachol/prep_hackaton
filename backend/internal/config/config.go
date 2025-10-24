package config

import "time"

type Config struct {
	ServerPort string
	DBUrl      string
	RedisAddr  string
	JWTSecret  string
	TokenTTL   time.Duration
}
