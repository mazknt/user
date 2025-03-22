package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	E "github.com/IBM/fp-go/either"

	email "authentication/Domain/models/User/Email"
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
		mockResponse   E.Either[error, dto.UserInformation]
		expectedStatus int
		expectedBody   dto.UserResponse
	}{
		{
			name:           "Success",
			request:        LOGIN_REQUEST["success"],
			mockResponse:   E.Right[error](SERVICE_RESPONSE["success"]),
			expectedStatus: http.StatusOK,
			expectedBody:   WANT["success"],
		},
		{
			name:           "Failure",
			request:        LOGIN_REQUEST["failed"],
			mockResponse:   E.Left[dto.UserInformation](errors.New("invalid auth code")),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   dto.UserResponse{Name: "", Email: "", Picture: "", Friends: make([]string, 0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("Login", tt.request.Code).Return(tt.mockResponse)

			requestBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mockController.Login(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			var actualResponse dto.UserResponse
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
		mockResponse   E.Either[error, dto.UserInformation]
		expectedStatus int
		expectedBody   dto.UserResponse
	}{
		{
			name:           "Success",
			request:        GET_USER_REQUEST["success"],
			mockResponse:   E.Right[error](SERVICE_RESPONSE["success"]),
			expectedStatus: http.StatusOK,
			expectedBody:   WANT["success"],
		},
		{
			name:           "BadRequest_MissingID",
			request:        GET_USER_REQUEST["failed"],
			mockResponse:   E.Left[dto.UserInformation](errors.New("missing ID")),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   dto.UserResponse{Name: "", Email: "", Picture: "", Friends: make([]string, 0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := E.Unwrap(email.NewEmail(tt.request.ID))
			if err == nil {
				mockService.On("GetUser", e).Return(tt.mockResponse)
			}

			requestBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/get-user-info", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mockController.GetUser(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			var actualResponse dto.UserResponse
			err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}
