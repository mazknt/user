package service_test

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/mock"
)

type MockFireStore struct {
	mock.Mock
}

func (m *MockFireStore) GetUserInformation(id email.Email) E.Either[error, user.User] {
	args := m.Called(id)
	return args.Get(0).(E.Either[error, user.User])
}

func (m *MockFireStore) SetUserInformation(u user.User) E.Either[error, user.User] {
	args := m.Called(u)
	return args.Get(0).(E.Either[error, user.User])
}
