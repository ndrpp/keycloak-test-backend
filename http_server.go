package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type httpServer struct {
	server *http.Server
}

func NewServer(host, port string, keycloak *keycloak) *httpServer {
	router := mux.NewRouter()

	noAuthRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Header.Get("Authorization") == ""
	}).Subrouter()

	authRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	controller := newController(keycloak)

	noAuthRouter.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res, _ := json.Marshal("Alive and well.")
		_, _ = w.Write(res)
	})
	noAuthRouter.HandleFunc("/login", func(writer http.ResponseWriter, req *http.Request) {
		controller.login(writer, req)
	})

	authRouter.HandleFunc("/docs", func(writer http.ResponseWriter, req *http.Request) {
		controller.getDocs(writer, req)
	})

	mdw := newMiddleware(keycloak)
	authRouter.Use(mdw.verifyToken)

	s := &httpServer{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      router,
			WriteTimeout: time.Hour,
			ReadTimeout:  time.Hour,
		},
	}

	return s
}

func (s *httpServer) listen() error {
	fmt.Println("Server is listening on: ", s.server.Addr)
	return s.server.ListenAndServe()
}
