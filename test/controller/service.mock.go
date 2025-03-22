package controller_test

import (
	email "authentication/Domain/models/User/Email"
	"authentication/dto"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Login(authCode string) E.Either[error, dto.UserInformation] {
	args := m.Called(authCode)
	return args.Get(0).(E.Either[error, dto.UserInformation])
}
func (m *MockService) GetUser(id email.Email) E.Either[error, dto.UserInformation] {
	args := m.Called(id)
	return args.Get(0).(E.Either[error, dto.UserInformation])
}

// func (m *MockService) GetUser(idE E.Either[error, string]) E.Either[error, dto.GetUserInfoResponse] {
// 	args := m.Called(idE)
// 	return E.Chain(
// 		func(id string) E.Either[error, dto.GetUserInfoResponse] {
// 			if id != "" {
// 				return E.Right[error](args.Get(0).(dto.GetUserInfoResponse))
// 			}
// 			return E.Left[dto.GetUserInfoResponse](errors.New("ログイン失敗"))
// 		})(idE)
// }
