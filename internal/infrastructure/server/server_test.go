package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/raphael-foliveira/go-table-tests/internal/adapters/api"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/dto"
	"github.com/raphael-foliveira/go-table-tests/internal/infrastructure/server"
	"github.com/raphael-foliveira/go-table-tests/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServer_ErrorHandler(t *testing.T) {
	app := server.CreateApp()
	apiGroup := app.Group("/api")
	mockUsersService := mocks.NewMockUsersService(t)
	usersHandler := api.NewUsersHandler(mockUsersService)

	usersHandler.SetupRoutes(apiGroup)

	req := httptest.NewRequest("POST", "/api/users/login", strings.NewReader(`{"invalid_field": "invalid_value"}`))
	recorder := httptest.NewRecorder()

	app.Server.Handler.ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	var responseBody *dto.ErrorResponse

	json.NewDecoder(response.Body).Decode(&responseBody)
	fmt.Printf("%+v\n", responseBody)
	assert.Equal(t, api.ErrInvalidPayload.Message, responseBody.Message)
	assert.Equal(t, http.StatusBadRequest, responseBody.StatusCode)
}
