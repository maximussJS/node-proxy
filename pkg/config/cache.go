package config

type ICacheConfig interface {
	UseRedis() bool
}

type CacheConfig struct {
	useRedis bool
}

func (cc *CacheConfig) UseRedis() bool {
	return cc.useRedis
}

var cacheConfig *CacheConfig

func init() {
	initCacheConfig()
}

func initCacheConfig() {
	useRedis := EnvOptionalBool("USE_REDIS", true)

	cacheConfig = &CacheConfig{
		useRedis: useRedis,
	}
}

func SingletonCacheConfig() *CacheConfig {
	if cacheConfig == nil {
		initCacheConfig()
	}

	return cacheConfig
}
