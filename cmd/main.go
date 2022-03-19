package main

import (
	echoV4 "github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase"
	"github.com/wellingtonlope/ticket-api/internal/infra/http"
	"github.com/wellingtonlope/ticket-api/internal/infra/http/echo"
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
	server := &echo.Server{
		Echo: echoV4.New(),
	}

	(&http.Http{
		UseCases:      useCases,
		Server:        server,
		Repositories:  repositories,
		Authenticator: authenticator,
	}).Start(1323)
}
