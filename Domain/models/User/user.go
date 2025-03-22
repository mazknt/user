package user

import (
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
)

type User struct {
	name    name.Name       `json:"name"`
	email   email.Email     `json:"email"`
	picture picture.Picture `json:"picture"`
	friends []email.Email   `json:"friends"`
}

func NewUser(name name.Name) func(e email.Email) func(picture picture.Picture) func(friends []email.Email) User {
	return func(e email.Email) func(picture picture.Picture) func(friends []email.Email) User {
		return func(picture picture.Picture) func(friends []email.Email) User {
			return func(friends []email.Email) User {
				return User{
					name:    name,
					email:   e,
					picture: picture,
					friends: friends,
				}
			}
		}
	}
}

func (u User) GetName() string {
	return u.name.Value()
}

func (u User) GetEmail() string {
	return u.email.Value()
}

func (u User) GetPicture() string {
	return u.picture.Value()
}
func (u User) GetFriends() []string {
	emails := make([]string, 0)
	for _, email := range u.friends {
		emails = append(emails, email.Value())
	}
	return emails
}
func (u User) GetNameObject() name.Name {
	return u.name
}
func (u User) GetEmailObject() email.Email {
	return u.email
}
func (u User) GetPictureObject() picture.Picture {
	return u.picture
}
func (u User) GetFriendsObject() []email.Email {
	emails := make([]email.Email, 0)
	for _, email := range u.friends {
		emails = append(emails, email)
	}
	return emails
}

func (u User) UpdateName(n name.Name) User {
	return NewUser(n)(u.email)(u.picture)(u.friends)
}
func (u User) UpdateEmail(e email.Email) User {
	return NewUser(u.name)(e)(u.picture)(u.friends)
}
func (u User) UpdatePicture(p picture.Picture) User {
	return NewUser(u.name)(u.email)(p)(u.friends)
}
