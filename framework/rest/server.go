package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/application/usecase"
	"github.com/wellingtonlope/ticket-api/framework/db/local"
	"github.com/wellingtonlope/ticket-api/framework/rest/handler"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil && os.Getenv("APP_ENV") != "production" {
		log.Fatalf("Error loading .env file: %v", err)
	}

	repositories := local.Repositories{}

	useCases, err := usecase.GetUseCases(&repositories, "123", time.Hour*24)
	if err != nil {
		log.Fatalf("Error during server initialization: %v", err)
	}

	e := echo.New()
	handler.InitHandlers(e, useCases)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
