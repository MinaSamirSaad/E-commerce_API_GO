package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/MinaSamirSaad/ecommerce/services/users"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix(shared.API_PREFIX).Subrouter()
	userHandler := users.NewHandler(s.db)
	userHandler.RegisterRoutes(subRouter)
	log.Println("listing on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
