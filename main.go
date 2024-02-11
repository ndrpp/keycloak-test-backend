package main

import (
	"fmt"
	"keycloak-go-backend/src"
	"keycloak-go-backend/src/services"
	"keycloak-go-backend/src/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	config := src.Config{Host: "localhost", Port: "8081"}
	logger := utils.NewLogger()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	keycloak := services.NewKeycloak()

	s := src.NewServer(config, logger, keycloak)
	src.Listen(s, logger)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
