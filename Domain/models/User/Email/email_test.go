package email

import (
	"strings"
	"testing"

	E "github.com/IBM/fp-go/either"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLeft bool
		wantErr  string
	}{
		{
			name:     "正常ケース: 有効なメールアドレス",
			input:    "test@example.com",
			wantLeft: false,
		},
		{
			name:     "正常ケース: 大文字を含むメールアドレス (小文字に正規化される)",
			input:    "TEST@EXAMPLE.COM",
			wantLeft: false,
		},
		{
			name:     "異常ケース: @が含まれていない",
			input:    "invalid-email",
			wantLeft: true,
			wantErr:  "mail: missing '@' or angle-addr",
		},
		{
			name:     "異常ケース: ドメインが無効",
			input:    "test@.com",
			wantLeft: true,
			wantErr:  "mail: no angle-addr",
		},
		{
			name:     "異常ケース: ローカルパートが無効",
			input:    "@example.com",
			wantLeft: true,
			wantErr:  "mail: no address",
		},
		{
			name:     "異常ケース: 空文字",
			input:    "",
			wantLeft: true,
			wantErr:  "mail: no address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewEmail(tt.input)

			if tt.wantLeft {
				if E.IsRight(result) {
					t.Errorf("expected Left but got Right")
				}
				if E.IsLeft(result) {
					E.Fold(
						func(err error) string {
							t.Errorf("expected error: %v, got: %v", tt.wantErr, err.Error())
							return ""
						},
						func(value Email) string { return "" },
					)(result)
				} else {
					t.Errorf("expected error but got nil")
				}
			} else {
				if E.IsLeft(result) {
					t.Errorf("expected Right but got Left")
				}
				// if email, ok := result.Right().(Email); ok {
				// 	expected := strings.ToLower(tt.input)
				// 	if email.Value() != expected {
				// 		t.Errorf("expected value: %v, got: %v", expected, email.Value())
				// 	}
				// }
				if E.IsRight(result) {
					E.Fold(
						func(err error) string {
							t.Errorf("expected error: %v, got: %v", tt.wantErr, err.Error())
							return ""
						},
						func(value Email) string {
							expected := strings.ToLower(tt.input)
							if value.Value() != expected {
								t.Errorf("expected value: %v, got: %v", expected, value.Value())
							}
							return ""
						},
					)(result)
				} else {
					t.Errorf("expected Email but got nil")
				}
			}
		})
	}
}
