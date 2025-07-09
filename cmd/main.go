package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	srv "urlMasker/internal/service"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "urlMasker",
		Usage: "Маскирует URL в файле и сохраняет результат в другой файл",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "in",
				Usage:    "путь к входному файлу",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "путь к выходному файлу (по умолчанию output.txt)",
				Value: "output.txt",
			},
			&cli.IntFlag{
				Name:  "timeout",
				Usage: "таймаут обработки (секунды, по умолчанию 2)",
				Value: 2,
			},
			&cli.StringFlag{
				Name:  "mode",
				Usage: "режим обработки: semaphore или workerpool (по умолчанию semaphore)",
				Value: "semaphore",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "уровень логирования: debug, info, warn, error (по умолчанию info)",
				Value: "info",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			inputPath := cmd.String("in")
			outputPath := cmd.String("out")
			timeout := cmd.Int("timeout")
			mode := cmd.String("mode")
			logLevel := cmd.String("log-level")

			var level slog.Level
			switch logLevel {
			case "debug":
				level = slog.LevelDebug
			case "info":
				level = slog.LevelInfo
			case "warn":
				level = slog.LevelWarn
			case "error":
				level = slog.LevelError
			default:
				level = slog.LevelInfo
			}
			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

			// Graceful shutdown: контекст отменяется по сигналу или по таймауту
			ctxTimeout, cancelTimeout := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancelTimeout()

			ctxSignal, cancelSignal := signal.NotifyContext(ctxTimeout, syscall.SIGINT, syscall.SIGTERM)
			defer cancelSignal()

			producer := srv.NewFileProducer(inputPath)
			presenter := srv.NewFilePresenter(outputPath)
			service := srv.NewService(producer, presenter, logger)

			var err error
			if mode == "workerpool" {
				err = service.RunWithWorkerPool(ctxSignal)
			} else {
				err = service.Run(ctxSignal)
			}

			if err != nil {
				return fmt.Errorf("ошибка: %v", err)
			}

			fmt.Println("Обработка завершена. Результат записан в:", outputPath)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
