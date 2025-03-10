package service_test

import (
	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/mock"
	googleOAuth "google.golang.org/api/oauth2/v2"
)

type MockGoogleAPI struct {
	mock.Mock
}

func (m *MockGoogleAPI) GetUserInfo(authCodeE E.Either[error, string]) E.Either[error, *googleOAuth.Userinfo] {
	args := m.Called(authCodeE)
	return args.Get(0).(E.Either[error, *googleOAuth.Userinfo])
}
