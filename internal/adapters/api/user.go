package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/dto"
	"github.com/raphael-foliveira/go-table-tests/internal/core/domain"
	"github.com/raphael-foliveira/go-table-tests/internal/core/ports"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
)

type UsersHandler struct {
	usersService ports.UsersService
}

func NewUsersHandler(usersService ports.UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (h *UsersHandler) SetupRoutes(group *echo.Group) {
	usersGroup := group.Group("/users")
	usersGroup.POST("/login", h.Login)
	usersGroup.POST("/signup", h.Signup)
}

func (h *UsersHandler) Login(ctx echo.Context) error {
	var loginPayload dto.LoginPayload
	if err := ctx.Bind(&loginPayload); err != nil {
		return ErrInvalidPayload
	}

	loginResponse, err := h.usersService.Login(loginPayload.Email, loginPayload.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return ErrInvalidCredentials
		}
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.LoginResponse{
		Email:    loginResponse.Email,
		Username: loginResponse.Username,
	})
}

func (h *UsersHandler) Signup(ctx echo.Context) error {
	var signupPayload dto.SignupPayload
	if err := ctx.Bind(&signupPayload); err != nil {
		return ErrInvalidPayload
	}

	signupResponse, err := h.usersService.Signup(&domain.SignupPayload{
		Email:    signupPayload.Email,
		Username: signupPayload.Username,
		Password: signupPayload.Password,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusCreated, &dto.SignupResponse{
		ID:       signupResponse.ID,
		Email:    signupResponse.Email,
		Username: signupResponse.Username,
	})
}

var (
	ErrInvalidPayload     = echo.NewHTTPError(400, "invalid payload")
	ErrInvalidCredentials = echo.NewHTTPError(401, "invalid credentials")
)
