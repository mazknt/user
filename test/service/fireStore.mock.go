package service_test

import (
	"authentication/dto"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/mock"
	googleOAuth "google.golang.org/api/oauth2/v2"
)

type MockFireStore struct {
	mock.Mock
}

func (m *MockFireStore) GetUserData(emailE E.Either[error, string]) E.Either[error, dto.LoginResponse] {
	args := m.Called(emailE)
	return args.Get(0).(E.Either[error, dto.LoginResponse])
}

func (m *MockFireStore) SetUserData(userInfo *googleOAuth.Userinfo) E.Either[error, dto.LoginResponse] {
	args := m.Called(userInfo)
	return args.Get(0).(E.Either[error, dto.LoginResponse])
}
