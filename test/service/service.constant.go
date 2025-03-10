package service_test

import "authentication/dto"

var LOGIN_RESPONSE = map[string]dto.LoginResponse{
	"get_user": {
		Name:    "get_user user",
		Email:   "get_user@example.com",
		Picture: "http://login",
		Friends: []string{"get_user1", "get_user2"},
	},
	"set_user": {
		Name:    "set_user user",
		Email:   "set_user@example.com",
		Picture: "http://login",
		Friends: []string{"set_user1", "set_user2"},
	},
	"failed": {
		Name:    "",
		Email:   "",
		Picture: "",
		Friends: nil,
	},
}

var GET_USER_RESPONSE = map[string]dto.GetUserInfoResponse{
	"success": {
		Name:    "login user",
		Email:   "login@example.com",
		Picture: "http://login",
		Friends: []string{"login1", "login2"},
	},
	"failed": {
		Name:    "",
		Email:   "",
		Picture: "",
		Friends: nil,
	},
}

var LOGIN_REQUEST = map[string]string{
	"success": "login@example.com",
	"failed":  "",
}
