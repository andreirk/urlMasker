package utils

import "sync"

type Task func()

type WorkerPool struct {
	workers int
	tasks   chan Task
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	return &WorkerPool{
		workers: workers,
		tasks:   make(chan Task),
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.wg.Add(1)
	go func() {
		wp.tasks <- func() {
			task()
			wp.wg.Done()
		}
	}()
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workers; i++ {
		go func() {
			for task := range wp.tasks {
				task()
			}
		}()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.tasks)
}
