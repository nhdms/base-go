package dbtool

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/spf13/viper"
	"time"
)

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (rc *RedisConfig) GetConnectionOptions() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Username: rc.Username,
		Password: rc.Password,
		DB:       rc.DB,
	}
}

func CreateRedisConnection(config *RedisConfig) (*redis.Client, error) {
	if config == nil {
		config = &RedisConfig{}
		sub := viper.Sub("redis")
		if sub == nil {
			return nil, fmt.Errorf("redis config not found")
		}

		err := sub.Unmarshal(config)
		if err != nil {
			return nil, err
		}
	}

	client := redis.NewClient(config.GetConnectionOptions())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %w", err)
	}

	logger.DefaultLogger.Infof("Connected to Redis %v:%v", config.Host, config.Port)
	return client, nil
}
