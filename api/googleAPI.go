package api

import (
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

type GoogleAPIInterface interface {
	GetUserInfo(authCodeE E.Either[error, string]) E.Either[error, *googleOAuth.Userinfo]
}

func (g GoogleAPI) GetUserInfo(authCodeE E.Either[error, string]) E.Either[error, *googleOAuth.Userinfo] {
	// OAuth2 configの作成
	oauth2Config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return FP.Pipe3(
		authCodeE,
		getToken(oauth2Config),
		getOAuth2Service(oauth2Config),
		getUserInfo,
	)
}

func getToken(oauth2Config oauth2.Config) func(authCodeE E.Either[error, string]) E.Either[error, *oauth2.Token] {
	return func(authCodeE E.Either[error, string]) E.Either[error, *oauth2.Token] {
		return E.Chain(func(authCode string) E.Either[error, *oauth2.Token] {
			token, err := oauth2Config.Exchange(context.Background(), authCode)
			if err != nil {
				return E.Left[*oauth2.Token](fmt.Errorf("failed to get token: %w", err))
			}
			return E.Right[error](token)
		})(authCodeE)
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

func getUserInfo(oauth2ServiceE E.Either[error, *googleOAuth.Service]) E.Either[error, *googleOAuth.Userinfo] {
	return E.Chain(
		func(oauth2Service *googleOAuth.Service) E.Either[error, *googleOAuth.Userinfo] {
			userInfo, err := oauth2Service.Userinfo.Get().Do()
			if err != nil {
				return E.Left[*googleOAuth.Userinfo](fmt.Errorf("failed to get userInfo: %v", err.Error()))
			}
			return E.Right[error](userInfo)
		})(oauth2ServiceE)
}
