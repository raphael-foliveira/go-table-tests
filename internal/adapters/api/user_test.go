package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/api"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/dto"
	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
	"github.com/raphael-foliveira/go-table-tests/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpTestServer() *echo.Echo {
	return echo.New()
}

func setUpDependencies(t *testing.T) (*api.UsersHandler, *mocks.MockUsersService) {
	mockUsersService := mocks.NewMockUsersService(t)
	usersHandler := api.NewUsersHandler(mockUsersService)
	return usersHandler, mockUsersService
}

func TestUserHandler_SetupRoutes(t *testing.T) {
	app := setUpTestServer()
	usersHandler, _ := setUpDependencies(t)

	apiGroup := app.Group("/api")

	usersHandler.SetupRoutes(apiGroup)

	tests := []struct {
		testName      string
		expectedRoute string
	}{
		{
			testName:      "Login Route",
			expectedRoute: "/api/users/login",
		},
		{
			testName:      "Signup Route",
			expectedRoute: "/api/users/signup",
		},
	}

	for i, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			routes := app.Routes()
			assert.Equal(t, tt.expectedRoute, routes[i].Path)
		})
	}
}

func TestUserHandlerLogin_Success(t *testing.T) {
	tests := []struct {
		mockLoginResponse *domain.LoginResponse
		testName          string
		loginJson         string
	}{
		{
			testName:  "Success",
			loginJson: `{"email": "valid@user.com", "password": "unhashedPassword"}`,
			mockLoginResponse: &domain.LoginResponse{
				Username: "validuser",
				Email:    "valid@user.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			testServer := setUpTestServer()
			usersHandler, mockUsersService := setUpDependencies(t)

			mockUsersService.EXPECT().
				Login(mock.Anything, mock.Anything).
				Return(tt.mockLoginResponse, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(tt.loginJson))
			req.Header.Add("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			ctx := testServer.NewContext(req, recorder)

			err := usersHandler.Login(ctx)
			assert.NoError(t, err)

			var responseBody *dto.LoginResponse
			json.NewDecoder(recorder.Body).Decode(&responseBody)
			assert.Equal(t, tt.mockLoginResponse.Email, responseBody.Email)
			assert.Equal(t, tt.mockLoginResponse.Username, responseBody.Username)
		})
	}
}

func TestUserHandlerLogin_InvalidBody(t *testing.T) {
	tests := []struct {
		mockLoginResponse *domain.LoginResponse
		testName          string
		loginJson         string
	}{
		{
			testName:  "Success",
			loginJson: `{"email": "valid@user.com", "password": "unhashedPassword"}`,
			mockLoginResponse: &domain.LoginResponse{
				Username: "validuser",
				Email:    "valid@user.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			testServer := setUpTestServer()
			usersHandler, mockUsersService := setUpDependencies(t)

			mockUsersService.EXPECT().
				Login(mock.Anything, mock.Anything).
				Return(tt.mockLoginResponse, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(tt.loginJson))
			req.Header.Add("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			ctx := testServer.NewContext(req, recorder)

			err := usersHandler.Login(ctx)
			assert.NoError(t, err)

			var responseBody *dto.LoginResponse
			json.NewDecoder(recorder.Body).Decode(&responseBody)
			assert.Equal(t, tt.mockLoginResponse.Email, responseBody.Email)
			assert.Equal(t, tt.mockLoginResponse.Username, responseBody.Username)
		})
	}
}

func TestUserHandlerLogin_InvalidCredentials(t *testing.T) {
	tests := []struct {
		testName  string
		loginJson string
	}{
		{
			testName:  "Invalid Credentials",
			loginJson: `{"email": "valid@user.com", "password": "invalidPassword"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			testServer := setUpTestServer()
			usersHandler, mockUsersService := setUpDependencies(t)

			mockUsersService.EXPECT().
				Login(mock.Anything, mock.Anything).
				Return(nil, service.ErrInvalidCredentials)

			req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(tt.loginJson))
			req.Header.Add("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			ctx := testServer.NewContext(req, recorder)

			err := usersHandler.Login(ctx)
			assert.Equal(t, api.ErrInvalidCredentials, err)
		})
	}
}

func TestUserHandler_Signup(t *testing.T) {
	tests := []struct {
		expectedError         error
		signupServiceError    error
		signupServiceResponse *domain.SignupResponse
		testName              string
		signupUsername        string
		signupEmail           string
		signupPassword        string
	}{
		{
			signupUsername: "testuser",
			signupEmail:    "test@user.com",
			signupPassword: "unhashedPassword",
			signupServiceResponse: &domain.SignupResponse{
				ID:       1,
				Username: "testuser",
				Email:    "test@user.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			testServer := setUpTestServer()
			usersHandler, mockUsersService := setUpDependencies(t)

			mockUsersService.EXPECT().Signup(&domain.SignupPayload{
				Username: tt.signupUsername,
				Email:    tt.signupEmail,
				Password: tt.signupPassword,
			}).Return(tt.signupServiceResponse, tt.signupServiceError)

			signupJson := fmt.Sprintf(
				`{"username": "%s", "email": "%s", "password": "%s"}`,
				tt.signupUsername,
				tt.signupEmail,
				tt.signupPassword,
			)

			req := httptest.NewRequest(http.MethodPost, "/api/users/signup", strings.NewReader(signupJson))
			req.Header.Add("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			ctx := testServer.NewContext(req, recorder)

			err := usersHandler.Signup(ctx)
			assert.Equal(t, tt.expectedError, err)

			if tt.expectedError == nil {
				var responseBody *dto.SignupResponse
				json.NewDecoder(recorder.Body).Decode(&responseBody)
				assert.Equal(t, tt.signupUsername, responseBody.Username)
				assert.Equal(t, tt.signupEmail, responseBody.Email)
				assert.True(t, responseBody.ID != 0)
			}
		})
	}
}
