package routes

import (
	"keycloak-go-backend/src/controllers"
	"keycloak-go-backend/src/middleware"
	"net/http"
)

func AddRoutes(
	mux *http.ServeMux,
	controller *controllers.Controller,
) {
	mux.Handle("/login", controller.Login())
	mux.Handle("/docs", middleware.VerifyToken(controller.GetDocs()))
	mux.Handle("/healthz", controller.HandleHealthZ())
	mux.Handle("/", http.NotFoundHandler())
}
