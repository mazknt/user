package picture

import (
	"testing"

	E "github.com/IBM/fp-go/either"
)

func TestNewPicture(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLeft bool
		wantErr  string
	}{
		{
			name:     "正常ケース: 有効なHTTPS URL",
			input:    "https://example.com/image.png",
			wantLeft: false,
		},
		{
			name:     "異常ケース: HTTP URL",
			input:    "http://example.com/image.png",
			wantLeft: true,
			wantErr:  "only HTTPS URLs are allowed",
		},
		{
			name:     "異常ケース: 無効なURL (スキームなし)",
			input:    "example.com/image.png",
			wantLeft: true,
			wantErr:  "only HTTPS URLs are allowed",
		},
		{
			name:     "異常ケース: 空文字",
			input:    "https://example.com/image.png ",
			wantLeft: true,
			wantErr:  "blank is not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPicture(tt.input)

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
						func(picture Picture) string { return "" },
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
							return ""
						},
						func(picture Picture) string {
							if picture.Value() != tt.input {
								t.Errorf("expected value: %v, got: %v", tt.input, picture.Value())
							}
							return ""
						},
					)(result)
				} else {
					t.Errorf("expected Picture but got nil")
				}
			}
		})
	}
}
