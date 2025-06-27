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

	type task struct {
		idx  int
		text string
	}
	type result struct {
		idx  int
		text string
	}

	inCh := make(chan task)
	outCh := make(chan result)
	semaphore := utils.NewSemaphore(numOfGorutines)

	// Fan-in
	results := make([]string, len(data))
	var faninWg sync.WaitGroup
	faninWg.Add(1)
	go func() {
		defer faninWg.Done()
		for res := range outCh {
			results[res.idx] = res.text
		}
	}()

	// Workers
	var workerWg sync.WaitGroup
	for i := 0; i < numOfGorutines; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for t := range inCh {
				semaphore.Acquire()
				masked := s.mask(t.text)
				outCh <- result{t.idx, masked}
				semaphore.Release()
			}
		}()
	}

	// Producer: send tasks
	go func() {
		for idx, line := range data {
			inCh <- task{idx, line}
		}
		close(inCh)
	}()

	workerWg.Wait()
	close(outCh)
	faninWg.Wait()

	if err := s.pres.Present(results); err != nil {
		return fmt.Errorf("presenter error: %w", err)
	}
	return nil
}

func (s *Service) RunWithWorkerPool() error {
	data, err := s.prod.Produce()
	if err != nil {
		return fmt.Errorf("producer error: %w", err)
	}

	type task struct {
		idx  int
		text string
	}
	type result struct {
		idx  int
		text string
	}

	inCh := make(chan task)
	outCh := make(chan result)
	pool := utils.NewWorkerPool(numOfGorutines)
	pool.Run()

	// Fan-in
	results := make([]string, len(data))
	var faninWg sync.WaitGroup
	faninWg.Add(1)
	go func() {
		defer faninWg.Done()
		for res := range outCh {
			results[res.idx] = res.text
		}
	}()

	// Workers
	var workerWg sync.WaitGroup
	for i := 0; i < numOfGorutines; i++ {
		workerWg.Add(1)
		pool.AddTask(func() {
			defer workerWg.Done()
			for t := range inCh {
				masked := s.mask(t.text)
				outCh <- result{t.idx, masked}
			}
		})
	}

	// Producer: sent tasks
	go func() {
		for idx, line := range data {
			inCh <- task{idx, line}
		}
		close(inCh)
	}()

	workerWg.Wait()
	close(outCh)
	faninWg.Wait()

	if err := s.pres.Present(results); err != nil {
		return fmt.Errorf("presenter error: %w", err)
	}
	return nil
}

func (s *Service) Run() error {
	return s.RunWithSemaphore()
	// return s.RunWithWorkerPool()
}
