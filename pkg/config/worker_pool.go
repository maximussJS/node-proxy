package config

type IWorkerPoolConfig interface {
	GetMaxWorkers() int
}

type WorkerPoolConfig struct {
	maxWorkers int
}

func (wpc *WorkerPoolConfig) GetMaxWorkers() int {
	return wpc.maxWorkers
}

var workerPoolConfig *WorkerPoolConfig

func init() {
	initWorkerPoolConfig()
}

func initWorkerPoolConfig() {
	maxWorkers := EnvOptionalInt("WORKER_POOL_MAX_WORKERS", 1000)

	workerPoolConfig = &WorkerPoolConfig{
		maxWorkers: maxWorkers,
	}
}

func SingletonWorkerPoolConfig() *WorkerPoolConfig {
	if workerPoolConfig == nil {
		initWorkerPoolConfig()
	}

	return workerPoolConfig
}
