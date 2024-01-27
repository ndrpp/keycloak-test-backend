package main

import (
	"keycloak-go-backend/src"
	"keycloak-go-backend/src/services"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	keycloak := services.NewKeycloak()
	s := src.NewServer("localhost", "8081", keycloak)
	s.Listen()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
