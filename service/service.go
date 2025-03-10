package service

import (
	"authentication/api"
	"authentication/dto"
	"authentication/util/convert"

	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
	"google.golang.org/api/oauth2/v2"
)

type Service struct {
	GoogleAPI api.GoogleAPIInterface
	FireStore api.FireStoreInterface
}

type ServiceInterface interface {
	Login(authCodeE E.Either[error, string]) E.Either[error, dto.LoginResponse]
	GetUser(idE E.Either[error, string]) E.Either[error, dto.GetUserInfoResponse]
}

func NewService(GoogleAPI api.GoogleAPI, FireStore api.FireStore) Service {
	return Service{
		GoogleAPI: GoogleAPI,
		FireStore: FireStore,
	}
}

func (s Service) Login(authCodeE E.Either[error, string]) E.Either[error, dto.LoginResponse] {
	userInfoE := s.GoogleAPI.GetUserInfo(authCodeE)
	userFromDBE := s.getUserFromDB(userInfoE)
	return s.setUserIfNeed(userInfoE, userFromDBE)
}

func (s Service) GetUser(idE E.Either[error, string]) E.Either[error, dto.GetUserInfoResponse] {
	return FP.Pipe1(
		s.FireStore.GetUserData(idE),
		E.Map[error](convert.GetUserInfoResponseFromLoginResponse),
	)
}

func (s Service) getUserFromDB(userInfoE E.Either[error, *oauth2.Userinfo]) E.Either[error, dto.LoginResponse] {
	return FP.Pipe2(
		userInfoE,
		E.Map[error](func(userInfo *oauth2.Userinfo) string {
			return userInfo.Email
		}),
		s.FireStore.GetUserData,
	)
}

func (s Service) setUserIfNeed(userInfoE E.Either[error, *oauth2.Userinfo], userFromDBE E.Either[error, dto.LoginResponse]) E.Either[error, dto.LoginResponse] {
	return E.Chain(func(userInfo *oauth2.Userinfo) E.Either[error, dto.LoginResponse] {
		return E.Fold(
			func(err error) E.Either[error, dto.LoginResponse] {
				if err.Error() == "user is not exist" {
					return s.FireStore.SetUserData(userInfo)
				}
				return E.Left[dto.LoginResponse](err)
			},
			func(userFromDB dto.LoginResponse) E.Either[error, dto.LoginResponse] {
				return E.Right[error](userFromDB)
			},
		)(userFromDBE)
	})(userInfoE)
}
