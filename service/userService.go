package service

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	domainservice "authentication/Domain/models/User/domainService"
	"authentication/dto"
	"authentication/service/apiIF"
	"authentication/service/repositoryIF"

	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
)

type User struct {
	authorizeAPI   apiIF.AuthorizationServiceIF
	userRepository repositoryIF.UserIF
}

type UserServiceInterface interface {
	Login(authCode string) E.Either[error, dto.UserInformation]
	GetUser(id email.Email) E.Either[error, dto.UserInformation]
}

func NewService(authorizeAPI apiIF.AuthorizationServiceIF, userRepository repositoryIF.UserIF) User {
	return User{
		authorizeAPI:   authorizeAPI,
		userRepository: userRepository,
	}
}

func (s User) Login(authCode string) E.Either[error, dto.UserInformation] {
	userE := s.authorizeAPI.GetUserInfo(authCode)
	idE := E.Map[error](
		func(user user.User) email.Email {
			return user.GetEmailObject()
		})(userE)

	return FP.Pipe4(
		E.Right[error](domainservice.CheckDuplicate(s.userRepository)),
		E.Ap[E.Either[error, bool]](idE),
		E.Flatten,
		E.Chain(
			func(duplicate bool) E.Either[error, user.User] {
				if duplicate {
					return E.Chain(s.userRepository.GetUserInformation)(idE)
				}
				return E.Chain(s.userRepository.SetUserInformation)(userE)
			}),
		E.Map[error](dto.NewUserInformaiton),
	)
}

func (s User) GetUser(id email.Email) E.Either[error, dto.UserInformation] {
	return FP.Pipe1(
		s.userRepository.GetUserInformation(id),
		E.Map[error](dto.NewUserInformaiton),
	)
}
