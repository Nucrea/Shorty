package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type Config struct {
	AppUrl      string
	AppPort     uint16
	PostgresUrl string
	RedisUrl    string
}

func NewConfig() (*Config, error) {
	pgUrl := os.Getenv("SHORTY_POSTGRES_URL")
	if pgUrl == "" {
		return nil, fmt.Errorf("empty db url")
	}

	redisUrl := os.Getenv("SHORTY_REDIS_URL")
	if redisUrl == "" {
		return nil, fmt.Errorf("empty redis url")
	}

	baseUrl := os.Getenv("SHORTY_APP_URL")
	if baseUrl == "" {
		return nil, fmt.Errorf("empty app url")
	}

	appPortEnv := os.Getenv("SHORTY_APP_PORT")
	if appPortEnv == "" {
		return nil, fmt.Errorf("empty app port")
	}
	appPort, err := strconv.Atoi(appPortEnv)
	if err != nil {
		return nil, fmt.Errorf("error parsing app port")
	}
	if appPort > math.MaxUint16 {
		return nil, fmt.Errorf("app port is out of range")
	}

	return &Config{
		AppUrl:      baseUrl,
		AppPort:     uint16(appPort),
		PostgresUrl: pgUrl,
		RedisUrl:    redisUrl,
	}, nil
}
