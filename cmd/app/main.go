package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/raphael-foliveira/go-table-tests/internal/adapters/api"
	postgresRepository "github.com/raphael-foliveira/go-table-tests/internal/adapters/repository/postgres"
	"github.com/raphael-foliveira/go-table-tests/internal/adapters/security"
	"github.com/raphael-foliveira/go-table-tests/internal/core/service"
	"github.com/raphael-foliveira/go-table-tests/internal/infrastructure/config"
	"github.com/raphael-foliveira/go-table-tests/internal/infrastructure/database"
	"github.com/raphael-foliveira/go-table-tests/internal/infrastructure/server"
)

func main() {
	cfg, err := config.NewConfig(".env")
	if err != nil {
		panic(err)
	}

	app := server.CreateApp()

	apiGroup := app.Group("/api")

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	usersRepository := postgresRepository.NewUsers(db)
	hasher := security.NewHashingService()
	usersService := service.NewUsersService(usersRepository, hasher)
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
