package service

import (
	"fmt"
	"sync"
	"urlMasker/internal/utils"
)

const numOfGorutines = 10

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{prod: prod, pres: pres}
}

func (s *Service) mask(str string) string {
	var result []byte
	i := 0
	for i < len(str) {
		if str[i] == ' ' {
			result = append(result, ' ')
			i++
			continue
		}
		start := i
		for i < len(str) && str[i] != ' ' {
			i++
		}
		word := str[start:i]
		if len(word) >= 7 && string(word[:7]) == "http://" {
			result = append(result, []byte("http://")...)
			for j := 7; j < len(word); j++ {
				result = append(result, '*')
			}
		} else if len(word) >= 8 && string(word[:8]) == "https://" {
			result = append(result, []byte("http://")...)
			for j := 8; j < len(word); j++ {
				result = append(result, '*')
			}
		} else {
			result = append(result, word...)
		}
	}
	return string(result)
}

func (s *Service) RunWithSemaphore() error {
	data, err := s.prod.Produce()
	if err != nil {
		return fmt.Errorf("producer error: %w", err)
	}
	var masked = make([]string, len(data))
	var wg sync.WaitGroup
	semaphore := utils.NewSemaphore(numOfGorutines)

	for i, line := range data {
		wg.Add(1)
		semaphore.Acquire()
		go func(idx int, l string) {
			defer wg.Done()
			defer semaphore.Release()
			fmt.Printf("Launched worker %d\n", idx)
			masked[idx] = s.mask(l)
		}(i, line)
	}
	wg.Wait()
	if err := s.pres.Present(masked); err != nil {
		return fmt.Errorf("presenter error: %w", err)
	}
	return nil
}

func (s *Service) RunWithWorkerPool() error {
	data, err := s.prod.Produce()
	if err != nil {
		return fmt.Errorf("producer error: %w", err)
	}
	var masked = make([]string, len(data))
	pool := utils.NewWorkerPool(numOfGorutines)
	pool.Run()
	for i, line := range data {
		idx, l := i, line
		pool.AddTask(func() {
			fmt.Printf("Launched worker %d\n", idx)
			masked[idx] = s.mask(l)
		})
	}
	pool.Wait()
	if err := s.pres.Present(masked); err != nil {
		return fmt.Errorf("presenter error: %w", err)
	}
	return nil
}

func (s *Service) Run() error {
	return s.RunWithSemaphore()
	// return s.RunWithWorkerPool()
}
