package name

import (
	"errors"
	"regexp"

	E "github.com/IBM/fp-go/either"
)

type Name struct {
	value string
}

func NewName(value string) E.Either[error, Name] {
	if validateMin(value) {
		return E.Left[Name](errors.New("名前には2文字以上入力して下さい"))
	}
	if validateMax(value) {
		return E.Left[Name](errors.New("名前は50文字以内にして下さい"))
	}
	if checkCharacter(value) {
		return E.Left[Name](errors.New("許可されていない文字が入力されています"))
	}
	return E.Right[error](Name{value: value})
}

func validateMin(value string) bool {
	if len(value) < 2 {
		return true
	}
	return false
}

func validateMax(value string) bool {
	if len(value) > 50 {
		return true
	}
	return false
}

func checkCharacter(value string) bool {
	if !regexp.MustCompile(`^[a-zA-Z0-9_.\- ]+$`).MatchString(value) {
		return true
	}
	return false
}

func (n Name) Value() string {
	return n.value
}
