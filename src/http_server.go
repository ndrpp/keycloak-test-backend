package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"keycloak-go-backend/src/controllers"
	"keycloak-go-backend/src/middleware"
	"keycloak-go-backend/src/services"
	"net/http"
	"time"
)

type HttpServer struct {
	Server *http.Server
}

func NewServer(host, port string, keycloak *services.Keycloak) *HttpServer {
	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	noAuthRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Header.Get("Authorization") == ""
	}).Subrouter()

	authRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	controller := controllers.NewController(keycloak)

	noAuthRouter.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res, _ := json.Marshal("Alive and well.")
		_, _ = w.Write(res)
	})
	noAuthRouter.HandleFunc("/login", func(writer http.ResponseWriter, req *http.Request) {
		controller.Login(writer, req)
	})

	authRouter.HandleFunc("/docs", func(writer http.ResponseWriter, req *http.Request) {
		controller.GetDocs(writer, req)
	})

	mdw := middleware.NewMiddleware(keycloak)
	authRouter.Use(mdw.VerifyToken)

	s := &HttpServer{
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      router,
			WriteTimeout: time.Hour,
			ReadTimeout:  time.Hour,
		},
	}

	return s
}

func (s *HttpServer) Listen() error {
	fmt.Println("Server is listening on: ", s.Server.Addr)
	return s.Server.ListenAndServe()
}
