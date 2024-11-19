package config

type cacheConfig struct {
	UseRedis bool
}

var CacheConfig *cacheConfig

func init() {
	useRedis := EnvOptionalBool("USE_REDIS", true)

	CacheConfig = &cacheConfig{
		UseRedis: useRedis,
	}
}
