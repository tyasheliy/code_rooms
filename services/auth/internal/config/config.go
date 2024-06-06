package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	AppConfig
	StorageConfig
	LogConfig
	CacheConfig
}

type AppConfig struct {
	Build      string
	WebApiPort string
	TokenConfig
}

type TokenConfig struct {
	Secret     string
	Expiration time.Duration
}

type StorageConfig struct {
	Driver       string
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	MigrationDir string
}

type CacheConfig struct {
	Cleanup    time.Duration
	Expiration time.Duration
}

type LogConfig struct {
	Format string
	Level  string
}

func New() (*Config, error) {
	build, ok := os.LookupEnv("BUILD")

	if !ok || build == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}

		build = os.Getenv("BUILD")
	}

	vp := viper.New()

	vp.AutomaticEnv()
	vp.SetConfigName(build)
	vp.AddConfigPath(".")
	vp.AddConfigPath("configs")
	vp.AddConfigPath("/etc/auth")

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppConfig: AppConfig{
			Build:      build,
			WebApiPort: vp.GetString("app.webapi_port"),
			TokenConfig: TokenConfig{
				Secret:     vp.GetString("app_token_secret"),
				Expiration: vp.GetDuration("app.token.expiration") * time.Minute,
			},
		},
		StorageConfig: StorageConfig{
			Driver:       vp.GetString("storage_driver"),
			Host:         vp.GetString("storage_host"),
			Port:         vp.GetString("storage_port"),
			User:         vp.GetString("storage_user"),
			Password:     vp.GetString("storage_password"),
			Database:     vp.GetString("storage_database"),
			MigrationDir: vp.GetString("storage_migration_dir"),
		},
		CacheConfig: CacheConfig{
			Cleanup:    vp.GetDuration("cache.cleanup") * time.Minute,
			Expiration: vp.GetDuration("cache.expiration") * time.Minute,
		},
		LogConfig: LogConfig{
			Format: vp.GetString("log.format"),
			Level:  vp.GetString("log.level"),
		},
	}

	return cfg, nil
}
