package service_test

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
	"authentication/dto"
	"log"

	E "github.com/IBM/fp-go/either"
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
	nameOb := name.NewName(n)
	emailOb := email.NewEmail(e)
	pictureOb := picture.NewPicture(p)
	friendsOb := make([]E.Either[error, email.Email], 0)
	for _, em := range fs {
		friendsOb = append(friendsOb, email.NewEmail(em))
	}
	userOb := user.NewUserFromEither(nameOb, emailOb, pictureOb, friendsOb)
	user, err := E.Unwrap(userOb)
	if err != nil {
		log.Fatal(err)
	}
	return dto.NewUserInformaiton(user)
}

func newUser(n string, e string, p string, fs []string) user.User {
	nameOb := name.NewName(n)
	emailOb := email.NewEmail(e)
	pictureOb := picture.NewPicture(p)
	friendsOb := make([]E.Either[error, email.Email], 0)
	for _, em := range fs {
		friendsOb = append(friendsOb, email.NewEmail(em))
	}
	userOb := user.NewUserFromEither(nameOb, emailOb, pictureOb, friendsOb)
	user, err := E.Unwrap(userOb)
	if err != nil {
		log.Fatal(err)
	}
	return user
}
