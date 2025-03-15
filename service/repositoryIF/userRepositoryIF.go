package repositoryIF

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"

	E "github.com/IBM/fp-go/either"
)

type UserIF interface {
	GetUserInformation(idE E.Either[error, email.Email]) E.Either[error, user.User]
	SetUserInformation(userE E.Either[error, user.User]) E.Either[error, user.User]
	// UpdateUserInformaiton(userE E.Either[error, user.User]) E.Either[error, user.User]
	// DeleteUserInformaiton(userE E.Either[error, user.User]) E.Either[error, user.User]
}
