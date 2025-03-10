package json_util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	E "github.com/IBM/fp-go/either"
)

// モックデータ構造体
type MockRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var testData = MockRequest{Name: "Bob", Email: "bob@example.com"}

// ReadRequestの正常系テスト
func TestReadRequest_Success(t *testing.T) {
	body, _ := json.Marshal(testData)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	result := ReadRequest[MockRequest](req)

	E.Fold(
		func(err error) string {
			t.Errorf("expected success but got error: %v", err)
			return ""
		},
		func(req MockRequest) string {
			if req.Name != "Alice" || req.Email != "alice@example.com" {
				t.Errorf("unexpected result: got %v, want %v", req, testData)
			}
			return ""
		},
	)(result)
}

// ReadRequestの無効なJSONテスト
func TestReadRequest_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	result := ReadRequest[MockRequest](req)

	E.Fold(
		func(err error) string {
			if err.Error() != "invalid character 'i' looking for beginning of object key string" {
				t.Fatal("エラーメッセージが異なる")
			}
			return ""
		},
		func(req MockRequest) string {
			t.Fatal("expected error but got success")
			return ""
		},
	)(result)
}

// ReadRequestの空ボディテスト
func TestReadRequest_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	result := ReadRequest[MockRequest](req)

	E.Fold(
		func(err error) string {
			if err.Error() != "unexpected end of JSON input" {
				t.Fatal("エラーメッセージが異なる")
			}
			return ""
		},
		func(req MockRequest) string {
			t.Fatal("expected error but got success")
			return ""
		},
	)(result)
}

// JsonUnmarshalの正常系テスト
func TestJsonUnmarshal_Success(t *testing.T) {
	body, _ := json.Marshal(testData)
	result := JsonUnmarshal[MockRequest](body)
	E.Fold(
		func(err error) string {
			t.Fatal("expected success but got error")
			return ""
		},
		func(req MockRequest) string {
			assert.Equal(t, testData, req)
			return ""
		},
	)(result)
}

// JsonUnmarshalの無効なJSONテスト
func TestJsonUnmarshal_InvalidJSON(t *testing.T) {
	result := JsonUnmarshal[MockRequest]([]byte("{invalid json}"))
	E.Fold(
		func(err error) string {
			assert.Equal(t, "invalid character 'i' looking for beginning of object key string", err.Error())
			return ""
		},
		func(req MockRequest) string {
			t.Fatal("expected error but got success")
			return ""
		},
	)(result)
}

// JsonUnmarshalの空ボディテスト
func TestJsonUnmarshal_EmptyBody(t *testing.T) {
	result := JsonUnmarshal[MockRequest]([]byte{})
	E.Fold(
		func(err error) string {
			assert.Equal(t, "unexpected end of JSON input", err.Error())
			return ""
		},
		func(req MockRequest) string {
			t.Fatal("expected error but got success")
			return ""
		},
	)(result)
}
