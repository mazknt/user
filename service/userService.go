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
	Login(authCodeE E.Either[error, string]) E.Either[error, dto.UserInformation]
	GetUser(idE E.Either[error, email.Email]) E.Either[error, dto.UserInformation]
}

func NewService(authorizeAPI apiIF.AuthorizationServiceIF, userRepository repositoryIF.UserIF) User {
	return User{
		authorizeAPI:   authorizeAPI,
		userRepository: userRepository,
	}
}

func (s User) Login(authCodeE E.Either[error, string]) E.Either[error, dto.UserInformation] {
	userE := s.authorizeAPI.GetUserInfo(authCodeE)
	idE := E.Map[error](
		func(user user.User) email.Email {
			return user.GetEmailObject()
		})(userE)

	return FP.Pipe2(
		domainservice.CheckDuplicate(s.userRepository)(idE),
		E.Chain(
			func(duplicate bool) E.Either[error, user.User] {
				if duplicate {
					return s.userRepository.GetUserInformation(idE)
				}
				return s.userRepository.SetUserInformation(userE)
			}),
		E.Map[error](
			func(user user.User) dto.UserInformation {
				return dto.NewUserInformaiton(user)
			}),
	)
}

func (s User) GetUser(idE E.Either[error, email.Email]) E.Either[error, dto.UserInformation] {
	return FP.Pipe1(
		s.userRepository.GetUserInformation(idE),
		E.Map[error](
			func(user user.User) dto.UserInformation {
				return dto.NewUserInformaiton(user)
			}),
	)
}
