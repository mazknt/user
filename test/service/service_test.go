package service_test

import (
	"authentication/dto"
	"authentication/service"
	"errors"
	"testing"

	user "authentication/Domain/models/User"
	email "authentication/Domain/models/User/Email"
	name "authentication/Domain/models/User/Name"
	picture "authentication/Domain/models/User/Picture"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name      string
		request   E.Either[error, string]
		googleAPI struct {
			getUserInfo struct {
				request  E.Either[error, string]
				response E.Either[error, user.User]
			}
		}
		firestore struct {
			getUser struct {
				request  E.Either[error, email.Email]
				response E.Either[error, user.User]
			}
			setUser struct {
				request  E.Either[error, user.User]
				response E.Either[error, user.User]
			}
		}
		want E.Either[error, dto.UserInformation]
	}{
		// 成功: 新規ユーザー
		{
			name:    "成功: 新規ユーザー",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}{
					request: E.Right[error]("valid_auth_code"),
					response: user.NewUserFromEither(name.NewName("get_user user"),
						email.NewEmail("get_user@example.com"),
						picture.NewPicture("https://login"),
						A.Map(
							func(e string) E.Either[error, email.Email] { return email.NewEmail(e) },
						)([]string{"get@user1.com", "get@user2.com"})),
					// response: E.Right[error](user.User{Email: "login@example.com"}),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}
				setUser struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}
			}{
				getUser: struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}{
					request:  email.NewEmail("get_user@example.com"),
					response: E.Left[user.User](errors.New("user is not exist")),
				},
				setUser: struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}{
					request: user.NewUserFromEither(name.NewName("get_user user"),
						email.NewEmail("get_user@example.com"),
						picture.NewPicture("https://login"),
						A.Map(
							func(e string) E.Either[error, email.Email] { return email.NewEmail(e) },
						)([]string{"get@user1.com", "get@user2.com"})),
					response: E.Right[error](RESPONSE["set_user"]),
				},
			},
			want: E.Right[error](WANT["set_user"]),
		},

		// 成功: 登録済みユーザー
		{
			name:    "成功: 登録済みユーザー",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}{
					request: E.Right[error]("valid_auth_code"),
					response: user.NewUserFromEither(name.NewName("get_user user"),
						email.NewEmail("get_user@example.com"),
						picture.NewPicture("https://login"),
						A.Map(
							func(e string) E.Either[error, email.Email] { return email.NewEmail(e) },
						)([]string{"get@user1.com", "get@user2.com"})),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}
				setUser struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}
			}{
				getUser: struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}{
					request:  email.NewEmail("get_user@example.com"),
					response: E.Right[error](RESPONSE["get_user"]),
				},
				setUser: struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}{
					request:  E.Left[user.User](errors.New("not called")),
					response: E.Left[user.User](errors.New("not called")),
				},
			},
			want: E.Right[error](WANT["get_user"]),
		},

		// 失敗: authCode の取得に失敗
		{
			name:    "失敗: authCode の取得に失敗",
			request: E.Left[string](errors.New("authCode が無効")),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}{
					request:  E.Left[string](errors.New("authCode が無効")),
					response: E.Left[user.User](errors.New("authCode が無効")),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}
				setUser struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}
			}{
				getUser: struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}{
					request:  E.Left[email.Email](errors.New("authCode が無効")),
					response: E.Left[user.User](errors.New("authCode が無効")),
				},
				setUser: struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}{
					request:  E.Left[user.User](errors.New("authCode が無効")),
					response: E.Left[user.User](errors.New("authCode が無効")),
				},
			},
			want: E.Left[dto.UserInformation](errors.New("authCode が無効")),
		},

		// 失敗: OAuthからのユーザー情報取得を失敗するパターン
		{
			name:    "失敗: OAuthからのユーザー情報取得を失敗するパターン",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}{
					request:  E.Right[error]("valid_auth_code"),
					response: E.Left[user.User](errors.New("OAuthユーザーの取得に失敗")),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}
				setUser struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}
			}{
				getUser: struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}{
					request:  E.Left[email.Email](errors.New("OAuthユーザーの取得に失敗")),
					response: E.Left[user.User](errors.New("OAuthユーザーの取得に失敗")),
				},
				setUser: struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}{
					request:  E.Left[user.User](errors.New("OAuthユーザーの取得に失敗")),
					response: E.Left[user.User](errors.New("OAuthユーザーの取得に失敗")),
				},
			},
			want: E.Left[dto.UserInformation](errors.New("OAuthユーザーの取得に失敗")),
		},

		// 失敗: DBからの取得を失敗するパターン
		{
			name:    "失敗: DBからの取得を失敗するパターン",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, user.User]
				}{
					request: E.Right[error]("valid_auth_code"),
					response: user.NewUserFromEither(name.NewName("get_user user"),
						email.NewEmail("get_user@example.com"),
						picture.NewPicture("https://login"),
						A.Map(
							func(e string) E.Either[error, email.Email] { return email.NewEmail(e) },
						)([]string{"get@user1.com", "get@user2.com"})),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}
				setUser struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}
			}{
				getUser: struct {
					request  E.Either[error, email.Email]
					response E.Either[error, user.User]
				}{
					request:  email.NewEmail("get_user@example.com"),
					response: E.Left[user.User](errors.New("DBユーザーの取得に失敗")),
				},
				setUser: struct {
					request  E.Either[error, user.User]
					response E.Either[error, user.User]
				}{
					request:  E.Left[user.User](errors.New("")),
					response: E.Left[user.User](errors.New("DBユーザーの取得に失敗")),
				},
			},
			want: E.Left[dto.UserInformation](errors.New("DBユーザーの取得に失敗")),
		},
	}

	for _, tt := range tests {
		mockGoogleAPI := new(MockGoogleAPI) // モックのインスタンスを作成
		mockFireStore := new(MockFireStore) // モックのインスタンスを作成
		mockService := service.NewService(mockGoogleAPI, mockFireStore)
		mockFireStore.On("GetUserInformation", tt.firestore.getUser.request).Return(tt.firestore.getUser.response)
		mockFireStore.On("SetUserInformation", tt.firestore.setUser.request).Return(tt.firestore.setUser.response)
		mockGoogleAPI.On("GetUserInfo", tt.googleAPI.getUserInfo.request).Return(tt.googleAPI.getUserInfo.response)
		got := mockService.Login(tt.request)
		assert.Equal(t, tt.want, got)
	}
}
