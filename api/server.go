package api

import (
	"net/http"

	"github.com/0xsj/rw-go/internal"

	"github.com/0xsj/rw-go/config"
	db "github.com/0xsj/rw-go/db/sqlc"
	"github.com/gorilla/mux"
)

type Server struct {
	config config.Config
	router *mux.Router
	store  db.Querier
	log    internal.Logger
}

func NewServer(config config.Config, store db.Querier, log internal.Logger) *Server {
	server := &Server{
		config: config,
		router: mux.NewRouter(),
		store:  store,
		log:    log,
	}
	return server
}

func (s *Server) MountHandlers() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", s.RegisterUser).Methods("POST")
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
