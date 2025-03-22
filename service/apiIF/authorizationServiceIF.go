package apiIF

import (
	user "authentication/Domain/models/User"

	E "github.com/IBM/fp-go/either"
)

type AuthorizationServiceIF interface {
	GetUserInfo(authCode string) E.Either[error, user.User]
}
