package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilePresenter_Present_Success(t *testing.T) {
	fname := "test_output.txt"
	defer os.Remove(fname)

	fp := NewFilePresenter(fname)
	lines := []string{"line1", "line2"}
	err := fp.Present(lines)
	assert.NoError(t, err)

	data, err := os.ReadFile(fname)
	assert.NoError(t, err)
	assert.Equal(t, "line1\nline2\n", string(data))
}

func TestFilePresenter_Present_Empty(t *testing.T) {
	fname := "test_output_empty.txt"
	defer os.Remove(fname)

	fp := NewFilePresenter(fname)
	err := fp.Present([]string{})
	assert.NoError(t, err)

	data, err := os.ReadFile(fname)
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func TestFilePresenter_Present_Error(t *testing.T) {
	// Попытка записать в недопустимый путь
	fp := NewFilePresenter("/invalid_path/test.txt")
	err := fp.Present([]string{"line"})
	assert.Error(t, err)
}
