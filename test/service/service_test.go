package service_test

import (
	"authentication/dto"
	"authentication/service"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/oauth2/v2"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name      string
		request   E.Either[error, string]
		googleAPI struct {
			getUserInfo struct {
				request  E.Either[error, string]
				response E.Either[error, *oauth2.Userinfo]
			}
		}
		firestore struct {
			getUser struct {
				request  E.Either[error, string]
				response E.Either[error, dto.LoginResponse]
			}
			setUser struct {
				request  *oauth2.Userinfo
				response E.Either[error, dto.LoginResponse]
			}
		}
		want E.Either[error, dto.LoginResponse]
	}{
		// 成功: 新規ユーザー
		{
			name:    "成功: 新規ユーザー",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}{
					request:  E.Right[error]("valid_auth_code"),
					response: E.Right[error](&oauth2.Userinfo{Email: "login@example.com"}),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}
				setUser struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}
			}{
				getUser: struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}{
					request:  E.Right[error]("login@example.com"),
					response: E.Left[dto.LoginResponse](errors.New("user is not exist")),
				},
				setUser: struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}{
					request:  &oauth2.Userinfo{Email: "login@example.com"},
					response: E.Right[error](LOGIN_RESPONSE["set_user"]),
				},
			},
			want: E.Right[error](LOGIN_RESPONSE["set_user"]),
		},

		// 成功: 登録済みユーザー
		{
			name:    "成功: 登録済みユーザー",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}{
					request:  E.Right[error]("valid_auth_code"),
					response: E.Right[error](&oauth2.Userinfo{Email: "existing_user@example.com"}),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}
				setUser struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}
			}{
				getUser: struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}{
					request:  E.Right[error]("existing_user@example.com"),
					response: E.Right[error](LOGIN_RESPONSE["get_user"]),
				},
				setUser: struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}{
					request:  nil,
					response: E.Left[dto.LoginResponse](errors.New("not called")),
				},
			},
			want: E.Right[error](LOGIN_RESPONSE["get_user"]),
		},

		// 失敗: authCode の取得に失敗
		{
			name:    "失敗: authCode の取得に失敗",
			request: E.Left[string](errors.New("authCode が無効")),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}{
					request:  E.Left[string](errors.New("authCode が無効")),
					response: E.Left[*oauth2.Userinfo](errors.New("authCode が無効")),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}
				setUser struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}
			}{
				getUser: struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}{
					request:  E.Left[string](errors.New("authCode が無効")),
					response: E.Left[dto.LoginResponse](errors.New("authCode が無効")),
				},
				setUser: struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}{
					request:  nil,
					response: E.Left[dto.LoginResponse](errors.New("authCode が無効")),
				},
			},
			want: E.Left[dto.LoginResponse](errors.New("authCode が無効")),
		},

		// 失敗: OAuthからのユーザー情報取得を失敗するパターン
		{
			name:    "失敗: OAuthからのユーザー情報取得を失敗するパターン",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}{
					request:  E.Right[error]("valid_auth_code"),
					response: E.Left[*oauth2.Userinfo](errors.New("OAuthユーザーの取得に失敗")),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}
				setUser struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}
			}{
				getUser: struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}{
					request:  E.Left[string](errors.New("OAuthユーザーの取得に失敗")),
					response: E.Left[dto.LoginResponse](errors.New("OAuthユーザーの取得に失敗")),
				},
				setUser: struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}{
					request:  nil,
					response: E.Left[dto.LoginResponse](errors.New("OAuthユーザーの取得に失敗")),
				},
			},
			want: E.Left[dto.LoginResponse](errors.New("OAuthユーザーの取得に失敗")),
		},

		// 失敗: DBからの取得を失敗するパターン
		{
			name:    "失敗: DBからの取得を失敗するパターン",
			request: E.Right[error]("valid_auth_code"),
			googleAPI: struct {
				getUserInfo struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}
			}{
				getUserInfo: struct {
					request  E.Either[error, string]
					response E.Either[error, *oauth2.Userinfo]
				}{
					request:  E.Right[error]("valid_auth_code"),
					response: E.Right[error](&oauth2.Userinfo{Email: "test@example.com"}),
				},
			},
			firestore: struct {
				getUser struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}
				setUser struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}
			}{
				getUser: struct {
					request  E.Either[error, string]
					response E.Either[error, dto.LoginResponse]
				}{
					request:  E.Right[error]("test@example.com"),
					response: E.Left[dto.LoginResponse](errors.New("DBユーザーの取得に失敗")),
				},
				setUser: struct {
					request  *oauth2.Userinfo
					response E.Either[error, dto.LoginResponse]
				}{
					request:  nil,
					response: E.Left[dto.LoginResponse](errors.New("DBユーザーの取得に失敗")),
				},
			},
			want: E.Left[dto.LoginResponse](errors.New("DBユーザーの取得に失敗")),
		},
	}

	for _, tt := range tests {
		mockGoogleAPI := new(MockGoogleAPI) // モックのインスタンスを作成
		mockFireStore := new(MockFireStore) // モックのインスタンスを作成
		mockService := &service.Service{GoogleAPI: mockGoogleAPI, FireStore: mockFireStore}
		mockFireStore.On("GetUserData", tt.firestore.getUser.request).Return(tt.firestore.getUser.response)
		mockFireStore.On("SetUserData", tt.firestore.setUser.request).Return(tt.firestore.setUser.response)
		mockGoogleAPI.On("GetUserInfo", tt.googleAPI.getUserInfo.request).Return(tt.googleAPI.getUserInfo.response)
		got := mockService.Login(tt.request)
		assert.Equal(t, tt.want, got)
	}
}
