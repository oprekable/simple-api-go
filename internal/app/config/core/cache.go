package core

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRing struct {
	Addrs        map[string]string `mapstructure:"addrs"`
	Password     string            `mapstructure:"password"`
	PoolTimeout  time.Duration     `mapstructure:"pool_timeout"`
	DB           int               `mapstructure:"db"`
	MaxRetries   int               `mapstructure:"max_retries"`
	PoolSize     int               `mapstructure:"pool_size"`
	DialTimeout  time.Duration     `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration     `mapstructure:"read_timeout"`
	WriteTimeout time.Duration     `mapstructure:"write_timeout"`
	IsEnabled    bool              `mapstructure:"is_enabled"`
}

func (cfg *RedisRing) Options() *redis.RingOptions {
	return &redis.RingOptions{
		Addrs:        cfg.Addrs,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  cfg.PoolTimeout,
	}
}
