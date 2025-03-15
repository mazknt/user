package dto

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
)

type UserInformation struct {
	name    name.Name
	email   email.Email
	picture picture.Picture
	friends []email.Email
}

func NewUserInformaiton(user user.User) UserInformation {
	return UserInformation{
		name:    user.GetNameObject(),
		email:   user.GetEmailObject(),
		picture: user.GetPictureObject(),
		friends: user.GetFriendsObject(),
	}
}
