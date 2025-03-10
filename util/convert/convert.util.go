package convert

import "authentication/dto"

func GetUserInfoResponseFromLoginResponse(userInfo dto.LoginResponse) dto.GetUserInfoResponse {
	return dto.GetUserInfoResponse{
		Name:    userInfo.Name,
		Email:   userInfo.Email,
		Picture: userInfo.Picture,
		Friends: userInfo.Friends,
	}
}
