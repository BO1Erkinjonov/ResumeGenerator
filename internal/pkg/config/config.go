package config

import (
	"github.com/spf13/cast"
	"os"
	"time"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout int
	}
	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	Token struct {
		Secret    string
		AccessTTL time.Duration
	}

	OTLPCollector struct {
		Host string
		Port string
	}
}

func Token() string {
	c := Config{}
	c.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")
	return c.Token.Secret
}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "resume_generate_service")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":9090")
	config.Context.Timeout = cast.ToInt(getEnv("CONTEXT_TIMEOUT", "5"))

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "resume-generator")
	config.Server.Port = getEnv("SERVER_PORT", ":9050")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "+_+diyor2005+_+")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "generator_resume")

	// token configuration
	config.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")

	return &config, nil
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
