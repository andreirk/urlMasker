package main

import (
	"flag"
	"fmt"
	"os"

	"urlMasker/internal/app"
)

func main() {
	var inputPath, outputPath string
	flag.StringVar(&inputPath, "input", "", "путь к входному файлу")
	flag.StringVar(&outputPath, "output", "output.txt", "путь к выходному файлу (по умолчанию output.txt)")
	flag.Parse()

	if inputPath == "" {
		fmt.Println("Необходимо указать путь к входному файлу через --input")
		os.Exit(1)
	}

	producer := app.NewFileProducer(inputPath)
	presenter := app.NewFilePresenter(outputPath)
	service := app.NewService(producer, presenter)

	if err := service.Run(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Обработка завершена. Результат записан в:", outputPath)
}
