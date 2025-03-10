package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	E "github.com/IBM/fp-go/either"

	"authentication/controller"
	"authentication/dto"

	"github.com/stretchr/testify/assert"
)

var mockService *MockService
var mockController *controller.Controller

func mockSetup() {
	mockService = new(MockService)
	mockController = &controller.Controller{Service: mockService}
}

func TestLogin(t *testing.T) {
	mockSetup()

	tests := []struct {
		name           string
		request        dto.LoginRequest
		mockResponse   E.Either[error, dto.LoginResponse]
		expectedStatus int
		expectedBody   dto.LoginResponse
	}{
		{
			name:           "Success",
			request:        LOGIN_REQUEST["success"],
			mockResponse:   E.Right[error](LOGIN_RESPONSE["success"]),
			expectedStatus: http.StatusOK,
			expectedBody:   LOGIN_RESPONSE["success"],
		},
		{
			name:           "Failure",
			request:        LOGIN_REQUEST["failed"],
			mockResponse:   E.Left[dto.LoginResponse](errors.New("invalid auth code")),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   dto.LoginResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("Login", E.Right[error](tt.request.Code)).Return(tt.mockResponse)

			requestBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mockController.Login(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			var actualResponse dto.LoginResponse
			err := json.NewDecoder(recorder.Body).Decode(&actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}

func TestGetUser(t *testing.T) {
	mockSetup()

	tests := []struct {
		name           string
		request        dto.GetUserInfoRequest
		mockResponse   E.Either[error, dto.GetUserInfoResponse]
		expectedStatus int
		expectedBody   dto.GetUserInfoResponse
	}{
		{
			name:           "Success",
			request:        GET_USER_REQUEST["success"],
			mockResponse:   E.Right[error](GET_USER_RESPONSE["success"]),
			expectedStatus: http.StatusOK,
			expectedBody:   GET_USER_RESPONSE["success"],
		},
		{
			name:           "BadRequest_MissingID",
			request:        GET_USER_REQUEST["failed"],
			mockResponse:   E.Left[dto.GetUserInfoResponse](errors.New("missing ID")),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   dto.GetUserInfoResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("GetUser", E.Right[error](tt.request.ID)).Return(tt.mockResponse)

			requestBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/get-user-info", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mockController.GetUser(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			var actualResponse dto.GetUserInfoResponse
			err := json.NewDecoder(recorder.Body).Decode(&actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}
