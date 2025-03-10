package json_util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	E "github.com/IBM/fp-go/either"
)

func ReadRequest[T any](r *http.Request) E.Either[error, T] {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return E.Left[T](fmt.Errorf("failed to read request: %w", err))
	}
	return JsonUnmarshal[T](bodyBytes)
}

func JsonUnmarshal[T any](bodyBytes []byte) E.Either[error, T] {
	var req T
	err := json.Unmarshal(bodyBytes, &req)
	if err != nil {
		return E.Left[T](fmt.Errorf("failed to unmarshal json: %w", err))
	}
	return E.Right[error](req)
}
