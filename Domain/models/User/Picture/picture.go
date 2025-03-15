package picture

import (
	"errors"
	"strings"

	E "github.com/IBM/fp-go/either"
)

// Picture 値オブジェクト
type Picture struct {
	url string
}

func NewPicture(urlValue string) E.Either[error, Picture] {
	if validateStartWithHttps(urlValue) {
		return E.Left[Picture](errors.New("only HTTPS URLs are allowed"))
	}
	if validateBlank(urlValue) {
		return E.Left[Picture](errors.New("blank is not allowed"))
	}
	if validatePath(urlValue) {
		return E.Left[Picture](errors.New("need path"))
	}
	return E.Right[error](Picture{url: urlValue})
}

func validateStartWithHttps(url string) bool {
	if url[:8] != "https://" {
		return true
	}
	return false
}

func validateBlank(url string) bool {
	if strings.Contains(url, " ") {
		return true
	}
	return false
}

func validatePath(url string) bool {
	if len(url) == 8 {
		return true
	}
	return false
}

func (p Picture) Value() string {
	return p.url
}
