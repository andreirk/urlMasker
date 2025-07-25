package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

type MockPresenter struct {
	mock.Mock
}

func (m *MockPresenter) Present(lines []string) error {
	args := m.Called(lines)
	return args.Error(0)
}

func TestService_Run_Success(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	input := []string{"http://example.com", "no links here"}
	masked := []string{"http://***********", "no links here"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", masked).Return(nil)

	svc := NewService(prod, pres)
	err := svc.Run()
	assert.NoError(t, err)
	prod.AssertExpectations(t)
	pres.AssertExpectations(t)
}

func TestService_Run_ProducerError(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	prod.On("Produce").Return([]string(nil), errors.New("fail prod"))

	svc := NewService(prod, pres)
	err := svc.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "producer error")
	prod.AssertExpectations(t)
}

func TestService_Run_PresenterError(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	input := []string{"http://example.com"}
	masked := []string{"http://***********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", masked).Return(errors.New("fail pres"))

	svc := NewService(prod, pres)
	err := svc.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "presenter error")
	prod.AssertExpectations(t)
	pres.AssertExpectations(t)
}

func TestService_mask(t *testing.T) {
	svc := NewService(nil, nil)
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Test 1", "http://example.com", "http://***********"},
		{"Test 2", "visit http://test.com now", "visit http://******** now"},
		{"Test 3", "no links here", "no links here"},
		{"Test 4", "https://secure.com", "http://**********"},
		{"Test 5", "http and https http://a https://b", "http and https http://* http://*"},
		{"Test 6", "", ""},
		{"Test 7", "   ", "   "},
		{"Test 8", "http://", "http://"},
		{"Test 9", "https://", "http://"},
		{"Test 10", "http://a http://b", "http://* http://*"},
	}

	for _, tt := range tests {
		got := svc.mask(tt.input)
		assert.Equal(t, tt.expected, got, "input: %q", tt.input)
	}
}
