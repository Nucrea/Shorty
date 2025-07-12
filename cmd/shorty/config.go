package main

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppUrl      string
	AppPort     uint16
	ApiKey      string
	PostgresUrl string
	RedisUrl    string
	LogFile     string
	OTELUrl     string

	MinioEndpoint     string
	MinioAccessKey    string
	MinioAccessSecret string
}

func NewConfig(options ...ConfigOptions) (*Config, error) {
	opts := configOptions{}
	for _, option := range options {
		option.Apply(&opts)
	}

	getenv := os.Getenv
	if opts.envFilePath != nil {
		envs, err := parseEnvFile(*opts.envFilePath)
		if err != nil {
			return nil, err
		}
		getenv = func(key string) string {
			return envs[key]
		}
	}

	logFile := getenv("SHORTY_LOG_FILE")

	pgUrl := getenv("SHORTY_POSTGRES_URL")
	if pgUrl == "" {
		return nil, fmt.Errorf("empty postgres url")
	}
	if _, err := url.Parse(pgUrl); err != nil {
		return nil, fmt.Errorf("bad postgres url")
	}

	otelUrl := getenv("SHORTY_OPENTELEMETRY_URL")
	// if otelUrl == "" {
	// 	return nil, fmt.Errorf("empty otel url")
	// }
	// if _, err := url.Parse(otelUrl); err != nil {
	// 	return nil, fmt.Errorf("bad otel url")
	// }

	redisUrl := getenv("SHORTY_REDIS_URL")
	if redisUrl == "" {
		return nil, fmt.Errorf("empty redis url")
	}
	if _, err := url.Parse(redisUrl); err != nil {
		return nil, fmt.Errorf("bad redis url")
	}

	appUrl := getenv("SHORTY_APP_URL")
	if appUrl == "" {
		return nil, fmt.Errorf("empty app url")
	}
	if _, err := url.Parse(appUrl); err != nil {
		return nil, fmt.Errorf("bad app url")
	}

	apiKey := getenv("SHORTY_APP_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("empty api key")
	}

	minioEndpoint := getenv("SHORTY_MINIO_ENDPOINT")
	if minioEndpoint == "" {
		return nil, fmt.Errorf("empty minio endpoint")
	}

	minioAccessKey := getenv("SHORTY_MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		return nil, fmt.Errorf("empty minio access key")
	}

	minioAccessSecret := getenv("SHORTY_MINIO_ACCESS_SECRET")
	if minioAccessSecret == "" {
		return nil, fmt.Errorf("empty minio secret")
	}

	appPortEnv := getenv("SHORTY_APP_PORT")
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
		AppUrl:            appUrl,
		AppPort:           uint16(appPort),
		ApiKey:            apiKey,
		PostgresUrl:       pgUrl,
		RedisUrl:          redisUrl,
		LogFile:           logFile,
		OTELUrl:           otelUrl,
		MinioEndpoint:     minioEndpoint,
		MinioAccessKey:    minioAccessKey,
		MinioAccessSecret: minioAccessSecret,
	}, nil
}

func parseEnvFile(envFilePath string) (map[string]string, error) {
	envFile, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("err opening env file: %w", err)
	}

	envs := map[string]string{}
	lines := strings.Split(string(envFile), "\n")
	for _, line := range lines {
		keyValuePair := strings.SplitN(line, "=", 2)
		if len(keyValuePair) != 2 {
			return nil, fmt.Errorf("bad env file, line: %s", line)
		}

		key, value := keyValuePair[0], keyValuePair[1]
		for strings.HasPrefix(value, "\"") {
			value = strings.TrimPrefix(value, "\"")
		}
		for strings.HasSuffix(value, "\"") {
			value = strings.TrimSuffix(value, "\"")
		}

		envs[key] = value
	}

	return envs, nil
}

type configOptions struct {
	envFilePath *string
}

type ConfigOptions interface {
	Apply(options *configOptions)
}

type configOptionsWithEnvFile struct {
	envFilePath string
}

func (c *configOptionsWithEnvFile) Apply(options *configOptions) {
	options.envFilePath = &c.envFilePath
}

func ConfigWithEnvFile(envFilePath string) ConfigOptions {
	return &configOptionsWithEnvFile{envFilePath: envFilePath}
}
