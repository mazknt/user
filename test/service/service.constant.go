package service_test

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
	"authentication/dto"
	"log"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
)

var RESPONSE = map[string]user.User{
	"get_user": newUser("get_user user", "get_user@example.com", "https://login", []string{"get@user1.com", "get@user2.com"}),
	"set_user": newUser("set_user user", "set_user@example.com", "https://login", []string{"set@user1.com", "set@user2.com"}),
}
var WANT = map[string]dto.UserInformation{
	"get_user": newUserInformation("get_user user", "get_user@example.com", "https://login", []string{"get@user1.com", "get@user2.com"}),
	"set_user": newUserInformation("set_user user", "set_user@example.com", "https://login", []string{"set@user1.com", "set@user2.com"}),
}

// var GET_USER_RESPONSE = map[string]user.User{
// 	"success": {
// 		Name:    "login user",
// 		Email:   "login@example.com",
// 		Picture: "http://login",
// 		Friends: []string{"login1", "login2"},
// 	},
// 	"failed": {
// 		Name:    "",
// 		Email:   "",
// 		Picture: "",
// 		Friends: nil,
// 	},
// }

var LOGIN_REQUEST = map[string]string{
	"success": "login@example.com",
	"failed":  "",
}

func newUserInformation(n string, e string, p string, fs []string) dto.UserInformation {
	return F.Pipe1(
		newUser(n, e, p, fs),
		dto.NewUserInformaiton,
	)
}

func newUser(n string, e string, p string, fs []string) user.User {
	userE := F.Pipe4(
		E.Right[error](user.NewUser),
		E.Ap[func(e email.Email) func(picture picture.Picture) func(friends []email.Email) user.User](name.NewName(n)),
		E.Ap[func(picture picture.Picture) func(friends []email.Email) user.User](email.NewEmail(e)),
		E.Ap[func(friends []email.Email) user.User](picture.NewPicture(p)),
		E.Ap[user.User](
			F.Pipe2(
				fs,
				A.Map(email.NewEmail),
				E.SequenceArray,
			)),
	)
	user, err := E.Unwrap(userE)
	if err != nil {
		log.Fatal(err)
	}
	return user
}
