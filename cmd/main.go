package main

import (
	"flag"
	"fmt"
	"os"
)

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func splitBySpace(s string) []string {
	result := make([]string, 0)

	str := []byte(s)
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' {
			result = append(result, string(str[:i]))
			str = str[i+1:]
		}
	}
	result = append(result, string(str))
	return result
}

// repeat returns a string that is repeated n times
func repeat(s rune, n int) string {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = byte(s)
	}
	return string(result)
}

func join(result []string) string {
	resultStr := make([]byte, 0)
	for i, v := range result {
		resultStr = append(resultStr, []byte(v)...)
		if i < len(result)-1 {
			resultStr = append(resultStr, ' ')
		}
	}
	return string(resultStr)
}

func mask(str string) string {
	result := splitBySpace(str)

	fmt.Println(result)

	for i, v := range result {
		if hasPrefix(v, "http://") {
			length := len(result[i]) - 7
			fmt.Println(result[i], length)
			result[i] = fmt.Sprintf("http://%s", repeat('*', length))
		}
	}

	return join(result)
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
