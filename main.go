package main

import (
	"fmt"
	"keycloak-go-backend/src"
	"keycloak-go-backend/src/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	keycloak := services.NewKeycloak()
	s := src.NewServer("localhost", "8081", keycloak)
	src.Listen(s)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
