package config

import "time"

type IRedisConfig interface {
	GetUrl() string
	GetDefaultTTl() time.Duration
}

type RedisConfig struct {
	Url        string
	DefaultTTl int
}

func (rc *RedisConfig) GetUrl() string {
	return rc.Url
}

func (rc *RedisConfig) GetDefaultTTl() time.Duration {
	return time.Duration(rc.DefaultTTl) * time.Second
}

var redisConfig *RedisConfig

func init() {
	initRedisConfig()
}

func initRedisConfig() {
	url := EnvRequiredString("REDIS_URL")

	defaultTtl := EnvOptionalInt("REDIS_DEFAULT_TTL", 30)

	redisConfig = &RedisConfig{
		Url:        url,
		DefaultTTl: defaultTtl,
	}
}

func SingletonRedisConfig() *RedisConfig {
	if redisConfig == nil {
		initRedisConfig()
	}

	return redisConfig
}
