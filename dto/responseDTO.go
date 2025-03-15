package dto

type UserResponse struct {
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Picture string   `json:"picture"`
	Friends []string `json:"friends"`
}

func CreateUserResponse(userInfo UserInformation) UserResponse {
	friendslist := make([]string, 0)
	for _, friend := range userInfo.friends {
		friendslist = append(friendslist, friend.Value())
	}
	return UserResponse{
		Name:    userInfo.name.Value(),
		Email:   userInfo.email.Value(),
		Picture: userInfo.picture.Value(),
		Friends: friendslist,
	}
}
