package service

import (
	"os"
)

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
