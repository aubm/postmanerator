package tests

import "github.com/stretchr/testify/mock"

type MockThemesManager struct {
	mock.Mock
}

func (m *MockThemesManager) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
