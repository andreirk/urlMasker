package service

import (
	"context"
	"fmt"
	"log/slog"
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
	prod   Producer
	pres   Presenter
	logger *slog.Logger
}

func NewService(prod Producer, pres Presenter, logger *slog.Logger) *Service {
	return &Service{prod: prod, pres: pres, logger: logger}
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

func (s *Service) RunWithSemaphore(ctx context.Context) error {
	s.logger.Debug("Service started (semaphore mode)")

	data, err := s.prod.Produce()
	if err != nil {
		s.logger.Error("Producer error", slog.String("error", err.Error()))
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

	results := make([]string, len(data))
	var faninWg sync.WaitGroup
	faninWg.Add(1)
	go func() {
		s.logger.Debug("Fan-in goroutine started")
		defer faninWg.Done()
		for res := range outCh {
			results[res.idx] = res.text
		}
		s.logger.Debug("Fan-in goroutine finished")
	}()

	var workerWg sync.WaitGroup
	for i := 0; i < numOfGorutines; i++ {
		workerID := i
		workerWg.Add(1)
		go func(id int) {
			logger := s.logger.With(slog.Int("worker", id))
			logger.Debug("Worker started")
			defer func() {
				logger.Debug("Worker finished")
				workerWg.Done()
			}()
			for {
				select {
				case <-ctx.Done():
					logger.Warn("Worker canceled by context")
					return
				case t, ok := <-inCh:
					if !ok {
						logger.Debug("Input channel closed")
						return
					}
					semaphore.Acquire()
					logger.Debug("Processing task", slog.Int("task_idx", t.idx))
					masked := s.mask(t.text)
					outCh <- result{t.idx, masked}
					semaphore.Release()
				}
			}
		}(workerID)
	}

	go func() {
		s.logger.Debug("Producer goroutine started")
		for idx, line := range data {
			select {
			case <-ctx.Done():
				s.logger.Warn("Producer canceled by context")
				close(inCh)
				return
			default:
				inCh <- task{idx, line}
			}
		}
		close(inCh)
		s.logger.Debug("Producer goroutine finished")
	}()

	workerWg.Wait()
	close(outCh)
	faninWg.Wait()

	if ctx.Err() != nil {
		s.logger.Warn("Context canceled or timed out", slog.String("err", ctx.Err().Error()))
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	if err := s.pres.Present(results); err != nil {
		s.logger.Error("Presenter error", slog.String("error", err.Error()))
		return fmt.Errorf("presenter error: %w", err)
	}
	s.logger.Info("Service finished (semaphore mode)")
	return nil
}

func (s *Service) RunWithWorkerPool(ctx context.Context) error {
	s.logger.Debug("Service started (workerpool mode)")

	data, err := s.prod.Produce()
	if err != nil {
		s.logger.Error("Producer error", slog.String("error", err.Error()))
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

	results := make([]string, len(data))
	var faninWg sync.WaitGroup
	faninWg.Add(1)
	go func() {
		s.logger.Debug("Fan-in goroutine started")
		defer faninWg.Done()
		for res := range outCh {
			results[res.idx] = res.text
		}
		s.logger.Debug("Fan-in goroutine finished")
	}()

	var workerWg sync.WaitGroup
	for i := 0; i < numOfGorutines; i++ {
		workerID := i
		workerWg.Add(1)
		pool.AddTask(func(id int) func() {
			return func() {
				logger := s.logger.With(slog.Int("worker", id))
				logger.Debug("Worker started")
				defer func() {
					logger.Debug("Worker finished")
					workerWg.Done()
				}()
				for {
					select {
					case <-ctx.Done():
						logger.Warn("Worker canceled by context")
						return
					case t, ok := <-inCh:
						if !ok {
							logger.Debug("Input channel closed")
							return
						}
						logger.Debug("Processing task", slog.Int("task_idx", t.idx))
						masked := s.mask(t.text)
						outCh <- result{t.idx, masked}
					}
				}
			}
		}(workerID))
	}

	go func() {
		s.logger.Debug("Producer goroutine started")
		for idx, line := range data {
			select {
			case <-ctx.Done():
				s.logger.Warn("Producer canceled by context")
				close(inCh)
				return
			default:
				inCh <- task{idx, line}
			}
		}
		close(inCh)
		s.logger.Debug("Producer goroutine finished")
	}()

	workerWg.Wait()
	close(outCh)
	faninWg.Wait()

	if ctx.Err() != nil {
		s.logger.Warn("Context canceled or timed out", slog.String("err", ctx.Err().Error()))
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	if err := s.pres.Present(results); err != nil {
		s.logger.Error("Presenter error", slog.String("error", err.Error()))
		return fmt.Errorf("presenter error: %w", err)
	}
	s.logger.Info("Service finished (workerpool mode)")
	return nil
}

func (s *Service) Run(ctx context.Context) error {
	return s.RunWithSemaphore(ctx)
	// return s.RunWithWorkerPool()
}
