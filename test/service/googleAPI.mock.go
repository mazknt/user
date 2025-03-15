package service_test

import (
	user "authentication/Domain/models/User"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/mock"
)

type MockGoogleAPI struct {
	mock.Mock
}

func (m *MockGoogleAPI) GetUserInfo(authCodeE E.Either[error, string]) E.Either[error, user.User] {
	args := m.Called(authCodeE)
	return args.Get(0).(E.Either[error, user.User])
}
