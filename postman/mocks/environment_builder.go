package mocks_test

import (
	. "github.com/srgrn/postmanerator/postman"
	"github.com/stretchr/testify/mock"
)

type MockEnvironmentBuilder struct {
	mock.Mock
}

func (m *MockEnvironmentBuilder) FromFile(file string) (Environment, error) {
	args := m.Called(file)
	return args.Get(0).(Environment), args.Error(1)
}
