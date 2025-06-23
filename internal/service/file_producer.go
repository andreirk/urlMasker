package service

import (
	"bufio"
	"os"
)

type FileProducer struct {
	path string
}

func NewFileProducer(path string) *FileProducer {
	return &FileProducer{path: path}
}

func (fp *FileProducer) Produce() ([]string, error) {
	file, err := os.Open(fp.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if lines == nil {
		lines = []string{}
	}
	return lines, nil
}
