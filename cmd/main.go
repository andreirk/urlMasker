package main

import (
	"flag"
	"fmt"
	"os"
	s "strings"
)

func mask(str string) string {
	result := s.Split(str, " ")

	for i, v := range result {
		if s.HasPrefix(v, "http://") {
			result[i] = fmt.Sprintf("http://%s", s.Repeat("*", len(v)-7))
		}
		if s.HasPrefix(v, "https://") {
			result[i] = fmt.Sprintf("https://%s", s.Repeat("*", len(v)-8))
		}
	}

	return s.Join(result, " ")
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
