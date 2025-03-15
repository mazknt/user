package user

import (
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"

	E "github.com/IBM/fp-go/either"
)

type User struct {
	name    name.Name       `json:"name"`
	email   email.Email     `json:"email"`
	picture picture.Picture `json:"picture"`
	friends []email.Email   `json:"friends"`
}

func NewUserFromEither(
	nameE E.Either[error, name.Name],
	emailE E.Either[error, email.Email],
	pictureE E.Either[error, picture.Picture],
	friendsEs []E.Either[error, email.Email]) E.Either[error, User] {
	friendsE := sequence(friendsEs)

	return E.Chain(
		func(friends []email.Email) E.Either[error, User] {
			return E.Chain(
				func(name name.Name) E.Either[error, User] {
					return E.Chain(
						func(email email.Email) E.Either[error, User] {
							return E.Chain(
								func(picture picture.Picture) E.Either[error, User] {
									return E.Right[error](User{
										name:    name,
										email:   email,
										picture: picture,
										friends: friends,
									})
								},
							)(pictureE)
						},
					)(emailE)
				},
			)(nameE)
		},
	)(friendsE)
}

func NewUser(name name.Name, email email.Email, picture picture.Picture, friends []email.Email) User {
	return User{
		name:    name,
		email:   email,
		picture: picture,
		friends: friends,
	}
}

func (u User) GetName() string {
	return u.name.Value()
}
func GetNameEither(u E.Either[error, User]) E.Either[error, string] {
	return E.Map[error](
		func(user User) string {
			return user.GetName()
		})(u)
}

func (u User) GetEmail() string {
	return u.email.Value()
}
func GetEmailEither(u E.Either[error, User]) E.Either[error, string] {
	return E.Map[error](
		func(user User) string {
			return user.GetName()
		})(u)
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

func (u User) UpdateNameFromEither(nameE E.Either[error, name.Name]) E.Either[error, User] {
	return E.Map[error](
		func(name name.Name) User {
			return NewUser(name, u.email, u.picture, u.friends)
		})(nameE)
}
func (u User) UpdateEmailFromEither(emailE E.Either[error, email.Email]) E.Either[error, User] {
	return E.Map[error](
		func(email email.Email) User {
			return NewUser(u.name, email, u.picture, u.friends)
		})(emailE)
}
func (u User) UpdatePictureFromEither(pictureE E.Either[error, picture.Picture]) E.Either[error, User] {
	return E.Map[error](
		func(picture picture.Picture) User {
			return NewUser(u.name, u.email, picture, u.friends)
		})(pictureE)
}

func sequence[L any, R any](arr []E.Either[L, R]) E.Either[L, []R] {
	results := make([]R, 0, len(arr))

	for _, e := range arr {
		val, err := E.Unwrap(e)
		if E.IsLeft(e) {
			return E.Left[[]R](err)
		}
		results = append(results, val)
	}

	return E.Right[L, []R](results)
}
