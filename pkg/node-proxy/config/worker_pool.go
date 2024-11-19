package config

type workerPoolConfig struct {
	MaxWorkers int
}

var WorkerPoolConfig *workerPoolConfig

func init() {
	maxWorkers := EnvOptionalInt("WORKER_POOL_MAX_WORKERS", 1000)

	WorkerPoolConfig = &workerPoolConfig{
		MaxWorkers: maxWorkers,
	}
}
