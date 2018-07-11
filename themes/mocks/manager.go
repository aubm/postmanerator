package mocks_test

import (
	. "github.com/aubm/postmanerator/themes"
	"github.com/stretchr/testify/mock"
)

type MockThemeManager struct {
	mock.Mock
}

func (m *MockThemeManager) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockThemeManager) Download(themeName string) error {
	return m.Called(themeName).Error(0)
}

func (m *MockThemeManager) Delete(theme string) error {
	return m.Called(theme).Error(0)
}

func (m *MockThemeManager) Open(themeName string) (*Theme, error) {
	args := m.Called(themeName)
	return args.Get(0).(*Theme), args.Error(1)
}
