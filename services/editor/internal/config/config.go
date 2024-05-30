package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	env_file = ".env"

	DEV_BUILD   = "dev"
	DEBUG_BUILD = "debug"
	PROD_BUILD  = "prod"
)

type Config struct {
	Build string
	AppConfig
	CacheConfig
	LoggerConfig
}

type AppConfig struct {
	SourceDir  string
	SocketPort string
	WebApiPort string
}

type CacheConfig struct {
	Cleanup    time.Duration
	Expiration time.Duration
}

type LoggerConfig struct {
	Level  string
	Format string
}

func New() (*Config, error) {
	build, ok := os.LookupEnv("BUILD")
	if !ok {
		err := godotenv.Load(env_file)
		if err != nil {
			return nil, err
		}
		build = os.Getenv("BUILD")
	}

	vp := viper.New()

	vp.SetConfigName(build)
	vp.AddConfigPath(".")
	vp.AddConfigPath("./configs")

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Build: build,
		AppConfig: AppConfig{
			SourceDir:  vp.GetString("app.source_dir"),
			SocketPort: vp.GetString("app.socket_port"),
			WebApiPort: vp.GetString("app.webapi_port"),
		},
		CacheConfig: CacheConfig{
			Cleanup:    vp.GetDuration("cache.cleanup") * time.Minute,
			Expiration: vp.GetDuration("cache.expiration") * time.Minute,
		},
		LoggerConfig: LoggerConfig{
			Level:  vp.GetString("log.level"),
			Format: vp.GetString("log.format"),
		},
	}, nil
}
