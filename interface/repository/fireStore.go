package repository

import (
	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	FP "github.com/IBM/fp-go/function"
	"google.golang.org/api/option"
)

type FireStore struct {
}

func (f FireStore) GetUserData(emailE E.Either[error, email.Email]) E.Either[error, user.User] {
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

func (f FireStore) SetUserData(user user.User) E.Either[error, user.User] {
	clientE := newFireStoreClient()
	defer func() {
		E.Chain(func(client *firestore.Client) E.Either[error, string] {
			client.Close()
			return E.Right[error]("")
		})(clientE)
	}()

	return FP.Pipe1(
		clientE,
		setUserDocument(user),
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

func getUserDocument(emailE E.Either[error, email.Email]) func(E.Either[error, *firestore.Client]) E.Either[error, *firestore.DocumentSnapshot] {
	return func(clientE E.Either[error, *firestore.Client]) E.Either[error, *firestore.DocumentSnapshot] {
		return E.Chain(func(client *firestore.Client) E.Either[error, *firestore.DocumentSnapshot] {
			return E.Chain(func(email email.Email) E.Either[error, *firestore.DocumentSnapshot] {
				return getUserDocumentFromUserInfo(email, client)
			})(emailE)
		})(clientE)
	}
}

func getUserDocumentFromUserInfo(email email.Email, client *firestore.Client) E.Either[error, *firestore.DocumentSnapshot] {
	doc, err := client.Collection("user_collection").Doc(email.Value()).Get(context.Background())
	if err != nil {
		return E.Left[*firestore.DocumentSnapshot](fmt.Errorf("failed to get userDocument: %w", err))
	}
	return E.Right[error](doc)
}

func setUserDocument(userInfo user.User) func(E.Either[error, *firestore.Client]) E.Either[error, user.User] {
	return func(clientE E.Either[error, *firestore.Client]) E.Either[error, user.User] {
		return E.Chain(func(client *firestore.Client) E.Either[error, user.User] {
			docRef := client.Collection("user_collection").Doc(userInfo.GetEmail())
			_, err := docRef.Set(context.Background(), map[string]interface{}{
				"name":    userInfo.GetName(),
				"picture": userInfo.GetPicture(),
				"email":   userInfo.GetEmail(),
				"friends": []string{},
			})
			if err != nil {
				return E.Left[user.User](fmt.Errorf("failed to set userDocument: %w", err))
			}
			return E.Right[error](userInfo)
		})(clientE)
	}
}

func createGetUserResponse(docE E.Either[error, *firestore.DocumentSnapshot]) E.Either[error, user.User] {
	return E.Chain(
		func(doc *firestore.DocumentSnapshot) E.Either[error, user.User] {
			if !doc.Exists() {
				return E.Left[user.User](errors.New("user is not exist"))
			}

			friendsE := FP.Pipe1(convertFromInterface[string](doc, "friends"),
				A.Map(
					func(id string) E.Either[error, email.Email] {
						return email.NewEmail(id)
					},
				))

			return user.NewUserFromEither(
				name.NewName(doc.Data()["name"].(string)),
				email.NewEmail(doc.Data()["email"].(string)),
				picture.NewPicture(doc.Data()["picture"].(string)),
				friendsE,
			)
		})(docE)
}
