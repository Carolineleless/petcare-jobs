package job

import (
	"context"
	"log"
)

type WorkerPool struct {
	queue   *Queue
	workers int
}

func NewWorkerPool(queue *Queue, workers int) *WorkerPool {
	return &WorkerPool{
		queue:   queue,
		workers: workers,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	for {
		select {
		case job := <-wp.queue.Dequeue():
			log.Printf("[worker %d] processing job %s", id, job.ID)
			ProcessJob(ctx, job)
		case <-ctx.Done():
			log.Printf("[worker %d] shutting down", id)
			return
		}
	}
}
