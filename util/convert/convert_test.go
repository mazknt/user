package convert

import (
	"authentication/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfoResponseFromLoginResponse(t *testing.T) {
	// テスト用の LoginResponse を作成
	userInfo := dto.LoginResponse{
		Name:    "Alice",
		Email:   "alice@example.com",
		Picture: "https://example.com/alice.jpg",
		Friends: []string{"Bob", "Charlie"},
	}

	// 関数を実行
	result := GetUserInfoResponseFromLoginResponse(userInfo)

	// assert を使って結果が期待される値と一致するかチェック
	assert.Equal(t, userInfo.Name, result.Name, "Name should match")
	assert.Equal(t, userInfo.Email, result.Email, "Email should match")
	assert.Equal(t, userInfo.Picture, result.Picture, "Picture should match")
	assert.ElementsMatch(t, userInfo.Friends, result.Friends, "Friends list should match")
}
