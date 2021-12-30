package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/application/usecase"
	"github.com/wellingtonlope/ticket-api/framework/db/mongodb"
	"github.com/wellingtonlope/ticket-api/framework/rest/handler"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil && os.Getenv("APP_ENV") != "production" {
		log.Fatalf("Error loading .env file: %v", err)
	}

	repositories := mongodb.Repositories{
		UriConnection: os.Getenv("MONGO_URI"),
		Database:      os.Getenv("MONGO_DATABASE"),
	}

	durationHour, err := strconv.Atoi(os.Getenv("TOKEN_DURATION_HOUR"))
	if err != nil {
		durationHour = 24
	}

	useCases, err := usecase.GetUseCases(&repositories, os.Getenv("APP_SECRET"), time.Hour*time.Duration(durationHour))
	if err != nil {
		log.Fatalf("Error during server initialization: %v", err)
	}

	e := echo.New()
	handler.InitHandlers(e, useCases)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
