package main

import (
	"flag"
	"fmt"
	"os"
)

func mask(str string) string {
	var result []byte
	i := 0
	for i < len(str) {
		// Пропускаем пробелы
		if str[i] == ' ' {
			result = append(result, ' ')
			i++
			continue
		}
		// Начало слова
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
		} else {
			result = append(result, word...)
		}
	}
	return string(result)
}

func main() {
	var input string

	flag.StringVar(&input, "input", "", "input string")
	flag.Parse()

	if input == "" {
		fmt.Println("input is required")
		os.Exit(1)
	}

	fmt.Println(mask(input))
}
