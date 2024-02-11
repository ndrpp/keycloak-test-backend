package routes

import (
	"keycloak-go-backend/src/controllers"
	"keycloak-go-backend/src/middleware"
	"keycloak-go-backend/src/utils"
	"net/http"
)

func AddRoutes(
	mux *http.ServeMux,
	logger *utils.Logger,
	userController *controllers.UserController,
) {
	mux.Handle("/login", userController.Login(logger))
	mux.Handle("/docs", middleware.VerifyToken(logger, userController.GetDocs(logger)))
	mux.Handle("/healthz", userController.HandleHealthZ(logger))
	mux.Handle("/", http.NotFoundHandler())
}
