package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// APIServer ...
type APIServer struct {
	config  *Config
	router  *mux.Router
	storage *sql.DB
}

// newAPIServer ...
func newAPIServer(config *Config) *APIServer {
	storage, err := openDB(config.Server.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &APIServer{
		config:  config,
		router:  mux.NewRouter(),
		storage: storage,
	}
}

// start ...
func (s *APIServer) start() error {

	s.configRouter()

	log.Printf("APIServer is starting: \n\tUsed config: %s\n\tListening on %s\n", configPath, s.config.Server.BindAddr)

	// Connecting to DB

	return http.ListenAndServe(s.config.Server.BindAddr, s.router)
}

// configRouter sets the serving paths and according handler functions
func (s *APIServer) configRouter() {
	s.router.HandleFunc("/", s.handleRoot())
}
