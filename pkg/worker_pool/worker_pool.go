package worker_pool

import (
	"github.com/alitto/pond/v2"
	"go.uber.org/dig"
	"json-rpc-node-proxy/pkg/config"
	"log"
)

type IWorkerPool[T any] interface {
	SubmitWithResult(fn func() (T, error)) (T, error)
	Submit(fn func() (T, error))
	Shutdown()
}

type WorkerPool[T any] struct {
	pool pond.ResultPool[T]
}

type WorkerPoolDependencies struct {
	dig.In

	Cfg config.IWorkerPoolConfig `name:"WorkerPoolConfig"`
}

func NewWorkerPool[T any](deps WorkerPoolDependencies) *WorkerPool[T] {
	maxWorkers := deps.Cfg.GetMaxWorkers()

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
