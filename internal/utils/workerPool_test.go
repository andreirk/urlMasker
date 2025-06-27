package utils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_Basic(t *testing.T) {
	var counter int32
	pool := NewWorkerPool(3)
	pool.Run()
	tasks := 10
	for i := 0; i < tasks; i++ {
		pool.AddTask(func() {
			atomic.AddInt32(&counter, 1)
		})
	}
	pool.Wait()
	assert.Equal(t, int32(tasks), counter)
}

func TestWorkerPool_Parallelism(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Run()
	start := time.Now()
	tasks := 4
	for i := 0; i < tasks; i++ {
		pool.AddTask(func() {
			time.Sleep(100 * time.Millisecond)
		})
	}
	pool.Wait()
	dur := time.Since(start)
	assert.GreaterOrEqual(t, dur.Milliseconds(), int64(200)) // минимум 2 волны по 2 задачи
}

func TestWorkerPool_ManyTasks(t *testing.T) {
	var counter int32
	pool := NewWorkerPool(5)
	pool.Run()
	tasks := 100
	for i := 0; i < tasks; i++ {
		pool.AddTask(func() {
			atomic.AddInt32(&counter, 1)
		})
	}
	pool.Wait()
	assert.Equal(t, int32(tasks), counter)
}
