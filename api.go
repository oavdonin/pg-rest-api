package main

import (
	"database/sql"
	"encoding/json"
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
	s.router.HandleFunc("/people/{uuid}", s.handleGetPerson()).Methods("GET")
}

// respond ...
func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		result, _ := json.Marshal(data)
		log.Println(string(result))
		w.Write(result)
	}
}

// error ...
func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// handleGetPerson bu UUID
func (s *APIServer) handleGetPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]
		person, err := s.FindByUUID(uuid)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, person)
	}
}
