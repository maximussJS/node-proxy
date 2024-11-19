package config

type redisConfig struct {
	Url        string
	DefaultTTl int
}

var RedisConfig *redisConfig

func init() {
	url := EnvRequiredString("REDIS_URL")

	defaultTtl := EnvOptionalInt("REDIS_DEFAULT_TTL", 30)

	RedisConfig = &redisConfig{
		Url:        url,
		DefaultTTl: defaultTtl,
	}
}
