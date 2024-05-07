package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/raphael-foliveira/go-table-tests/internal/adapters/api"
	"github.com/raphael-foliveira/go-table-tests/internal/application/server"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
)

func main() {
	app := server.CreateApp()

	apiGroup := app.Group("/api")

	usersService := &service.UsersService{}
	usersHandler := api.NewUsersHandler(usersService)

	usersHandler.SetupRoutes(apiGroup)

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	go func() {
		err := server.Start(app, 3000)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	<-ctx.Done()

	fmt.Println("server interrupted")
}
