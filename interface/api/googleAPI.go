package api

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
	"context"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth "google.golang.org/api/oauth2/v2"
)

type GoogleAPI struct{}

func (g GoogleAPI) GetUserInfo(authCode string) E.Either[error, user.User] {
	// OAuth2 configの作成
	oauth2Config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return FP.Pipe4(
		authCode,
		getToken(oauth2Config),
		getOAuth2Service(oauth2Config),
		E.Chain(getOAuthUserInfo),
		E.Chain(convertToUserInfo),
	)
}

func getToken(oauth2Config oauth2.Config) func(authCode string) E.Either[error, *oauth2.Token] {
	return func(authCode string) E.Either[error, *oauth2.Token] {
		token, err := oauth2Config.Exchange(context.Background(), authCode)
		if err != nil {
			return E.Left[*oauth2.Token](fmt.Errorf("failed to get token: %w", err))
		}
		return E.Right[error](token)
	}
}

func getOAuth2Service(oauth2Config oauth2.Config) func(tokenE E.Either[error, *oauth2.Token]) E.Either[error, *googleOAuth.Service] {
	return func(tokenE E.Either[error, *oauth2.Token]) E.Either[error, *googleOAuth.Service] {
		return E.Chain(
			func(token *oauth2.Token) E.Either[error, *googleOAuth.Service] {
				client := oauth2Config.Client(context.Background(), token)
				oauth2Service, err := googleOAuth.New(client)
				if err != nil {
					E.Left[*oauth2.Token](fmt.Errorf("failed to get oauth2Service: %v", err.Error()))
				}
				return E.Right[error](oauth2Service)
			})(tokenE)
	}
}

func getOAuthUserInfo(oauth2Service *googleOAuth.Service) E.Either[error, *googleOAuth.Userinfo] {
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return E.Left[*googleOAuth.Userinfo](fmt.Errorf("failed to get userInfo: %v", err.Error()))
	}
	return E.Right[error](userInfo)
}

func convertToUserInfo(oauthUser *googleOAuth.Userinfo) E.Either[error, user.User] {
	return FP.Pipe4(
		E.Right[error](user.NewUser),
		E.Ap[func(e email.Email) func(picture picture.Picture) func(friends []email.Email) user.User](name.NewName(oauthUser.Name)),
		E.Ap[func(picture picture.Picture) func(friends []email.Email) user.User](email.NewEmail(oauthUser.Email)),
		E.Ap[func(friends []email.Email) user.User](picture.NewPicture(oauthUser.Picture)),
		E.Ap[user.User](E.Right[error](make([]email.Email, 0))),
	)

}
