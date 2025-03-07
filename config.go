package main

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	AppUrl      string
	AppPort     uint16
	PostgresUrl string
	RedisUrl    string
	LogFile     string
}

func NewConfig() (*Config, error) {
	logFile := os.Getenv("SHORTY_LOG_FILE")

	pgUrl := os.Getenv("SHORTY_POSTGRES_URL")
	if pgUrl == "" {
		return nil, fmt.Errorf("empty postgres url")
	}
	if _, err := url.Parse(pgUrl); err != nil {
		return nil, fmt.Errorf("bad postgres url")
	}

	redisUrl := os.Getenv("SHORTY_REDIS_URL")
	if redisUrl == "" {
		return nil, fmt.Errorf("empty redis url")
	}
	if _, err := url.Parse(redisUrl); err != nil {
		return nil, fmt.Errorf("bad redis url")
	}

	appUrl := os.Getenv("SHORTY_APP_URL")
	if appUrl == "" {
		return nil, fmt.Errorf("empty app url")
	}
	if _, err := url.Parse(appUrl); err != nil {
		return nil, fmt.Errorf("bad app url")
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
		AppUrl:      appUrl,
		AppPort:     uint16(appPort),
		PostgresUrl: pgUrl,
		RedisUrl:    redisUrl,
		LogFile:     logFile,
	}, nil
}
