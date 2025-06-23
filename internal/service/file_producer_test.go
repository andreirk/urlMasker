package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileProducer_Produce_Success(t *testing.T) {
	fname := "test_input.txt"
	content := "line1\nline2\n"
	err := os.WriteFile(fname, []byte(content), 0644)
	assert.NoError(t, err)
	t.Cleanup(func() { os.Remove(fname) })

	fp := NewFileProducer(fname)
	lines, err := fp.Produce()
	assert.NoError(t, err)
	assert.Equal(t, []string{"line1", "line2"}, lines)
}

func TestFileProducer_Produce_FileNotFound(t *testing.T) {
	fp := NewFileProducer("not_exists.txt")
	lines, err := fp.Produce()
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestFileProducer_Produce_EmptyFile(t *testing.T) {
	fname := "empty.txt"
	err := os.WriteFile(fname, []byte(""), 0644)
	assert.NoError(t, err)
	t.Cleanup(func() { os.Remove(fname) })

	fp := NewFileProducer(fname)
	lines, err := fp.Produce()
	assert.NoError(t, err)
	assert.Equal(t, []string{}, lines)
}
