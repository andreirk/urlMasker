package utils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSemaphore_Basic(t *testing.T) {
	sem := NewSemaphore(2)
	var counter int32

	done := make(chan struct{})
	for i := 0; i < 5; i++ {
		go func() {
			sem.Acquire()
			defer sem.Release()
			atomic.AddInt32(&counter, 1)
			time.Sleep(50 * time.Millisecond)

			done <- struct{}{}
		}()
	}
	for i := 0; i < 5; i++ {
		<-done
	}
	assert.Equal(t, int32(5), counter)
}

func TestSemaphore_Parallelism(t *testing.T) {
	sem := NewSemaphore(2)
	var active int32
	var maxActive int32
	done := make(chan struct{})
	tasks := 6
	for i := 0; i < tasks; i++ {
		go func() {
			sem.Acquire()
			defer sem.Release()
			n := atomic.AddInt32(&active, 1)
			if n > atomic.LoadInt32(&maxActive) {
				atomic.StoreInt32(&maxActive, n)
			}
			time.Sleep(30 * time.Millisecond)
			atomic.AddInt32(&active, -1)

			done <- struct{}{}
		}()
	}
	for i := 0; i < tasks; i++ {
		<-done
	}
	assert.Equal(t, int32(2), maxActive)
}
