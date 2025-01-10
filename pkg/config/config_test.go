package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavissig/ConfigLib-GO/pkg/config"
)

func TestLoadConfig_LoadOneFile(t *testing.T) {
	type RedisConfig struct {
		Address  string `envconfig:"DB_REDIS_ADDRESS" required:"true"`
		Password string `envconfig:"DB_REDIS_PASSWORD" required:"true"`
	}

	expectedCfg := &RedisConfig{
		Address:  "redis.storage.address:6379",
		Password: "redisPass123Example",
	}

	t.Run("LoadOneFile", func(t *testing.T) {
		files := []string{
			"fixtures/.env.test",
		}

		cfg, err := config.LoadConfig[RedisConfig](
			config.AddFiles(files),
			config.WithPrintConfig(),
		)
		assert.NoError(t, err)
		assert.Equal(t, cfg, expectedCfg)
	})
}

func TestLoadConfig_LoadTwoFiles(t *testing.T) {
	type RedisConfig struct {
		Address  string `envconfig:"DB_REDIS_ADDRESS" required:"true"`
		Password string `envconfig:"DB_REDIS_PASSWORD" required:"true"`
	}

	type PostgresConfig struct {
		Address  string `envconfig:"DB_POSTGRES_ADDRESS" required:"true"`
		Password string `envconfig:"DB_POSTGRES_PASSWORD" required:"true"`
	}

	type DBConfig struct {
		Redis    RedisConfig
		Postgres PostgresConfig
	}

	expectedCfg := &DBConfig{
		Redis: RedisConfig{
			Address:  "redis.storage.address:6379",
			Password: "redisPass123Example",
		},
		Postgres: PostgresConfig{
			Address:  "psql.storage.address:6379",
			Password: "passPgSQL123Example",
		},
	}

	t.Run("LoadOneFile", func(t *testing.T) {
		files := []string{
			"fixtures/.env.test",
			"fixtures/.env.testTwo",
		}

		cfg, err := config.LoadConfig[DBConfig](
			config.AddFiles(files),
			config.WithPrintConfig(),
		)
		assert.NoError(t, err)
		assert.Equal(t, cfg, expectedCfg)
	})
}
