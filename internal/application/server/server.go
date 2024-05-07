package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/dto"
)

func CreateApp() *echo.Echo {
	app := echo.New()
	app.HTTPErrorHandler = customErrorHandler
	return app
}

func Start(app *echo.Echo, port int) error {
	addr := fmt.Sprintf(":%d", port)
	return app.Start(addr)
}

func customErrorHandler(err error, ctx echo.Context) {
	res := &dto.ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		res.StatusCode = httpErr.Code
		res.Message = httpErr.Message.(string)
	}

	ctx.JSON(res.StatusCode, res)
}
