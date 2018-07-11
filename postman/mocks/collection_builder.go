package mocks_test

import (
	. "github.com/aubm/postmanerator/postman"
	"github.com/stretchr/testify/mock"
)

type MockCollectionBuilder struct {
	mock.Mock
}

func (m *MockCollectionBuilder) FromFile(file string, options BuilderOptions) (Collection, error) {
	args := m.Called(file, options)
	return args.Get(0).(Collection), args.Error(1)
}
