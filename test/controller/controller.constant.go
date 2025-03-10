package controller_test

import "authentication/dto"

var LOGIN_RESPONSE = map[string]dto.LoginResponse{
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

var LOGIN_REQUEST = map[string]dto.LoginRequest{
	"success": {Code: "test123"},
	"failed":  {Code: ""},
}

var GET_USER_REQUEST = map[string]dto.GetUserInfoRequest{
	"success": {ID: "login@example.com"},
	"failed":  {ID: ""},
}
