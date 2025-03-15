package name

import (
	"testing"

	E "github.com/IBM/fp-go/either"
)

func TestNewName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLeft bool
		wantErr  string
	}{
		{
			name:     "正常ケース: 2文字の英数字",
			input:    "ab",
			wantLeft: false,
		},
		{
			name:     "正常ケース: 50文字の英数字",
			input:    "mnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			wantLeft: false,
		},
		{
			name:     "異常ケース: 1文字未満",
			input:    "a",
			wantLeft: true,
			wantErr:  "名前には2文字以上入力して下さい",
		},
		{
			name:     "異常ケース: 51文字以上",
			input:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890",
			wantLeft: true,
			wantErr:  "名前は50文字以内にして下さい",
		},
		{
			name:     "異常ケース: 許可されていない文字が含まれる",
			input:    "name@invalid!",
			wantLeft: true,
			wantErr:  "許可されていない文字が入力されています",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewName(tt.input)
			if tt.wantLeft {
				if E.IsRight(result) {
					t.Errorf("expected Left but got Right")
				}
				if E.IsLeft(result) {
					E.Fold(
						func(err error) string {
							if err.Error() != tt.wantErr {
								t.Errorf("expected error: %v, got: %v", tt.wantErr, err.Error())
							}
							return ""
						},
						func(name Name) string {
							return ""
						},
					)(result)
				} else {
					t.Errorf("expected error but got nil")
				}
			} else {
				if E.IsLeft(result) {
					t.Errorf("expected Right but got Left")
				}
				if E.IsRight(result) {
					E.Fold(
						func(err error) string {
							if err.Error() != tt.wantErr {
								t.Errorf("expected Name but got nil")
							}
							return ""
						},
						func(name Name) string {
							if name.Value() != tt.input {
								t.Errorf("expected value: %v, got: %v", tt.input, name.Value())
							}
							return ""
						},
					)(result)
				}
			}
		})
	}
}
