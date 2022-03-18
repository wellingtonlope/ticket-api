package main

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase"
	"github.com/wellingtonlope/ticket-api/internal/infra/http"
	"github.com/wellingtonlope/ticket-api/internal/infra/jwt"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
	"time"
)

func main() {
	repositories, err := (&memory.Repositories{}).GetRepositories()
	if err != nil {
		panic("Error during repositories connection")
	}

	useCases, err := usecase.GetUseCases(repositories)
	if err != nil {
		panic("Error during server initialization")
	}

	authenticator := jwt.NewAuthenticator(repositories.UserRepository, "secret", time.Hour*time.Duration(24))
	e := echo.New()
	http.Http{
		UseCases:      useCases,
		Echo:          e,
		Repositories:  repositories,
		Authenticator: authenticator,
	}.Init()

	e.Logger.Fatal(e.Start(":1323"))
}
