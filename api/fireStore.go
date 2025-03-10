package api

import (
	"authentication/dto"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
	googleOAuth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type FireStore struct {
}

type FireStoreInterface interface {
	GetUserData(emailE E.Either[error, string]) E.Either[error, dto.LoginResponse]
	SetUserData(userInfo *googleOAuth.Userinfo) E.Either[error, dto.LoginResponse]
}

func (f FireStore) GetUserData(emailE E.Either[error, string]) E.Either[error, dto.LoginResponse] {
	clientE := newFireStoreClient()
	defer func() {
		E.Chain(func(client *firestore.Client) E.Either[error, string] {
			client.Close()
			return E.Right[error]("")
		})(clientE)
	}()

	return FP.Pipe2(
		clientE,
		getUserDocument(emailE),
		createGetUserResponse,
	)
}

func (f FireStore) SetUserData(userInfo *googleOAuth.Userinfo) E.Either[error, dto.LoginResponse] {
	clientE := newFireStoreClient()
	defer func() {
		E.Chain(func(client *firestore.Client) E.Either[error, string] {
			client.Close()
			return E.Right[error]("")
		})(clientE)
	}()

	return FP.Pipe2(
		clientE,
		setUserDocument(userInfo),
		createSetUserResponse(userInfo),
	)
}

func newFireStoreClient() E.Either[error, *firestore.Client] {
	client, err := firestore.NewClient(context.Background(), os.Getenv("PROJECT_ID"), option.WithCredentialsFile(os.Getenv("CREDENTIAL_FILE")))
	if err != nil {
		return E.Left[*firestore.Client](fmt.Errorf("failed to create firestore client: %w", err))
	}
	return E.Right[error](client)
}

func convertFromInterface[T any](doc *firestore.DocumentSnapshot, key string) []T {
	// fruitsInterfaceは []interface{} 型です
	slice, ok := doc.Data()[key].([]interface{})
	if !ok {
		log.Fatalf("Error: not a slice of interfaces")
	}

	// []interface{} から []string に変換
	var strSlice []T
	for _, el := range slice {
		strEl, ok := el.(T)
		if !ok {
			log.Fatalf("Error: element %v is not a string", strEl)
		}
		strSlice = append(strSlice, strEl)
	}
	return strSlice
}

func getUserDocument(emailE E.Either[error, string]) func(E.Either[error, *firestore.Client]) E.Either[error, *firestore.DocumentSnapshot] {
	return func(clientE E.Either[error, *firestore.Client]) E.Either[error, *firestore.DocumentSnapshot] {
		return E.Chain(func(client *firestore.Client) E.Either[error, *firestore.DocumentSnapshot] {
			return E.Chain(func(email string) E.Either[error, *firestore.DocumentSnapshot] {
				return getUserDocumentFromUserInfo(email, client)
			})(emailE)
		})(clientE)
	}
}

func getUserDocumentFromUserInfo(email string, client *firestore.Client) E.Either[error, *firestore.DocumentSnapshot] {
	doc, err := client.Collection("user_collection").Doc(email).Get(context.Background())
	if err != nil {
		return E.Left[*firestore.DocumentSnapshot](fmt.Errorf("failed to get userDocument: %w", err))
	}
	return E.Right[error](doc)
}

func setUserDocument(userInfo *googleOAuth.Userinfo) func(E.Either[error, *firestore.Client]) E.Either[error, string] {
	return func(clientE E.Either[error, *firestore.Client]) E.Either[error, string] {
		return E.Chain(func(client *firestore.Client) E.Either[error, string] {
			docRef := client.Collection("user_collection").Doc(userInfo.Email)
			_, err := docRef.Set(context.Background(), map[string]interface{}{
				"name":    userInfo.Name,
				"picture": userInfo.Picture,
				"email":   userInfo.Email,
				"friends": []string{},
			})
			if err != nil {
				return E.Left[string](fmt.Errorf("failed to set userDocument: %w", err))
			}
			return E.Right[error]("ok")
		})(clientE)
	}
}

func createGetUserResponse(docE E.Either[error, *firestore.DocumentSnapshot]) E.Either[error, dto.LoginResponse] {
	return E.Chain(
		func(doc *firestore.DocumentSnapshot) E.Either[error, dto.LoginResponse] {
			if !doc.Exists() {
				return E.Left[dto.LoginResponse](errors.New("user is not exist"))
			}

			friends := convertFromInterface[string](doc, "friends")
			return E.Right[error](dto.LoginResponse{
				Name:    doc.Data()["name"].(string),
				Email:   doc.Data()["email"].(string),
				Picture: doc.Data()["picture"].(string),
				Friends: friends,
			})
		})(docE)
}

func createSetUserResponse(userInfo *googleOAuth.Userinfo) func(docE E.Either[error, string]) E.Either[error, dto.LoginResponse] {
	return func(docE E.Either[error, string]) E.Either[error, dto.LoginResponse] {
		return E.Chain(
			func(ok string) E.Either[error, dto.LoginResponse] {
				return E.Right[error](dto.LoginResponse{
					Name:    userInfo.Name,
					Email:   userInfo.Email,
					Picture: userInfo.Picture,
					Friends: []string{},
				})
			})(docE)
	}
}
