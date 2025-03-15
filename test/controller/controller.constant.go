package controller_test

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
	"authentication/dto"

	E "github.com/IBM/fp-go/either"
)

var SERVICE_RESPONSE = map[string]dto.UserInformation{
	"success": newUserInfoDTO(
		"login user",
		"login@example.com",
		"https://login",
		[]string{"login1@example.com", "login2@example.com"},
	),
}
var WANT = map[string]dto.UserResponse{
	"success": dto.CreateUserResponse(newUserInfoDTO(
		"login user",
		"login@example.com",
		"https://login",
		[]string{"login1@example.com", "login2@example.com"},
	)),
}

// var GET_USER_RESPONSE = map[string]dto.UserInformation{
// 	"success": {
// 		Name:    "login user",
// 		Email:   "login@example.com",
// 		Picture: "https://login",
// 		Friends: []string{"login1", "login2"},
// 	},
// 	"failed": {
// 		Name:    "",
// 		Email:   "",
// 		Picture: "",
// 		Friends: nil,
// 	},
// }

var LOGIN_REQUEST = map[string]dto.LoginRequest{
	"success": {Code: "test123"},
	"failed":  {Code: ""},
}

var GET_USER_REQUEST = map[string]dto.GetUserInfoRequest{
	"success": {ID: "login@example.com"},
	"failed":  {ID: ""},
}

func newUserInfoDTO(n string, e string, p string, friends []string) dto.UserInformation {
	nameE := name.NewName(n)
	emailE := email.NewEmail(e)
	pictureE := picture.NewPicture(p)
	friendsEs := make([]E.Either[error, email.Email], 0)
	for _, friend := range friends {
		friendsEs = append(friendsEs, email.NewEmail(friend))
	}
	userE := user.NewUserFromEither(nameE, emailE, pictureE, friendsEs)
	user, err := E.Unwrap(userE)
	if err != nil {
		panic(err)
	}
	return dto.NewUserInformaiton(user)
}
