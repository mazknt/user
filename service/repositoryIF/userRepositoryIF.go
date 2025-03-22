package repositoryIF

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"

	E "github.com/IBM/fp-go/either"
)

type UserIF interface {
	GetUserInformation(id email.Email) E.Either[error, user.User]
	SetUserInformation(user user.User) E.Either[error, user.User]
	// UpdateUserInformaiton(userE E.Either[error, user.User]) E.Either[error, user.User]
	// DeleteUserInformaiton(userE E.Either[error, user.User]) E.Either[error, user.User]
}
