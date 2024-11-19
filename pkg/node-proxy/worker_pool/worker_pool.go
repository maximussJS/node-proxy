package worker_pool

import (
	"github.com/alitto/pond/v2"
	"json-rpc-node-proxy/pkg/node-proxy/config"
	"log"
)

type WorkerPool[T any] struct {
	pool pond.ResultPool[T]
}

func NewWorkerPool[T any]() *WorkerPool[T] {
	maxWorkers := config.WorkerPoolConfig.MaxWorkers

	if maxWorkers <= 0 {
		log.Fatalf("Invalid maxWorkers value: %d", maxWorkers)
	}

	return &WorkerPool[T]{
		pool: pond.NewResultPool[T](maxWorkers),
	}
}

func (wp *WorkerPool[T]) SubmitWithResult(fn func() (T, error)) (T, error) {
	task := wp.pool.SubmitErr(fn)

	return task.Wait()
}

func (wp *WorkerPool[T]) Submit(fn func() (T, error)) {
	task := wp.pool.SubmitErr(fn)

	task.Wait()
}

func (wp *WorkerPool[T]) Shutdown() {
	wp.pool.StopAndWait()
}
