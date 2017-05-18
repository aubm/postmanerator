package tests

import "github.com/stretchr/testify/mock"

type MockThemesManager struct {
	mock.Mock
}

func (m *MockThemesManager) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockThemesManager) Download(themeName string) error {
	return m.Called(themeName).Error(0)
}

func (m *MockThemesManager) Delete(theme string) error {
	return m.Called(theme).Error(0)
}
