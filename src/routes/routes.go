package routes

import (
	"keycloak-go-backend/src/controllers"
	"keycloak-go-backend/src/middleware"
	"net/http"
)

func AddRoutes(
	mux *http.ServeMux,
	userController *controllers.UserController,
) {
	mux.Handle("/login", userController.Login())
	mux.Handle("/docs", middleware.VerifyToken(userController.GetDocs()))
	mux.Handle("/healthz", userController.HandleHealthZ())
	mux.Handle("/", http.NotFoundHandler())
}
