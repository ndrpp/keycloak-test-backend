package src

import (
	"fmt"
	"keycloak-go-backend/src/controllers"
	"keycloak-go-backend/src/middleware"
	"keycloak-go-backend/src/routes"
	"keycloak-go-backend/src/services"
	"keycloak-go-backend/src/utils"
	"net/http"
	"time"
)

type Config struct {
	Host string
	Port string
}

func NewServer(config Config, logger *utils.Logger, keycloak *services.Keycloak) *http.Server {
	mux := http.NewServeMux()
	userController := controllers.NewUserController(keycloak)
	routes.AddRoutes(mux, logger, userController)
	var handler http.Handler = mux
	handler = middleware.CorsMiddeware(handler)

	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler:      handler,
		WriteTimeout: time.Hour,
		ReadTimeout:  time.Hour,
	}

	return s
}

func Listen(s *http.Server, logger *utils.Logger) error {
	logger.Info(fmt.Sprintf("Server is listening on: %s", s.Addr))
	return s.ListenAndServe()
}
