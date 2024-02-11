package src

import (
	"fmt"
	"keycloak-go-backend/src/controllers"
	"keycloak-go-backend/src/middleware"
	"keycloak-go-backend/src/routes"
	"keycloak-go-backend/src/services"
	"net/http"
	"time"
)

func NewServer(host, port string, keycloak *services.Keycloak) *http.Server {
	mux := http.NewServeMux()
	controller := controllers.NewController(keycloak)
	routes.AddRoutes(mux, controller)
	var handler http.Handler = mux
	handler = middleware.CorsMiddeware(handler)

	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      handler,
		WriteTimeout: time.Hour,
		ReadTimeout:  time.Hour,
	}

	return s
}

func Listen(s *http.Server) error {
	fmt.Println("Server is listening on: ", s.Addr)
	return s.ListenAndServe()
}
