package app

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
	return lines, nil
}

type FilePresenter struct {
	path string
}

func NewFilePresenter(path string) *FilePresenter {
	return &FilePresenter{path: path}
}

func (fp *FilePresenter) Present(lines []string) error {
	file, err := os.Create(fp.path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
