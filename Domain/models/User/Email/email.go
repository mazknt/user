package email // Email 値オブジェクト
import (
	"net/mail"
	"strings"

	E "github.com/IBM/fp-go/either"
)

type Email struct {
	value string
}

func NewEmail(value string) E.Either[error, Email] {
	value = strings.ToLower(value)
	_, err := mail.ParseAddress(value)
	if err != nil {
		return E.Left[Email](err)
	}
	return E.Right[error](Email{value})
}

func (e Email) Value() string {
	return e.value
}
