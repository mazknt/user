package domainservice

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	"authentication/service/repositoryIF"

	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
)

func CheckDuplicate(userRepository repositoryIF.UserIF) func(idE E.Either[error, email.Email]) E.Either[error, bool] {
	return func(idE E.Either[error, email.Email]) E.Either[error, bool] {
		return FP.Pipe1(
			userRepository.GetUserInformation(idE),
			E.Fold(
				func(err error) E.Either[error, bool] {
					if err.Error() == "user is not exist" {
						return E.Right[error](false)
					}
					return E.Left[bool](err)
				},
				func(user user.User) E.Either[error, bool] { return E.Right[error](true) },
			),
		)
	}
}
