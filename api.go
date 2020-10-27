package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// responseWriter represents http server return code
type responseWriter struct {
	http.ResponseWriter
	code int
}

// APIServer structure
type APIServer struct {
	config  *Config
	router  *mux.Router
	storage *sql.DB
}

// newAPIServer init function
func newAPIServer(config *Config) *APIServer {
	storage, err := openDB(config.Server.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	CsvToDB("titanic.csv", storage)

	return &APIServer{
		config:  config,
		router:  mux.NewRouter(),
		storage: storage,
	}
}

// start APIServer function
func (s *APIServer) start() error {

	s.configRouter()

	log.Printf("APIServer is starting: \n\tUsed config: %s\n\tListening on %s\n", configPath, s.config.Server.BindAddr)

	// Connecting to DB

	return http.ListenAndServe(s.config.Server.BindAddr, s.router)
}

// configRouter sets the serving paths and according handler functions
func (s *APIServer) configRouter() {
	s.router.HandleFunc("/people/{uuid}", s.handleGetPerson()).Methods("GET")
	s.router.HandleFunc("/people", s.handleGetPeople()).Methods("GET")
	s.router.HandleFunc("/people", s.handleAddPerson()).Methods("POST")
	s.router.HandleFunc("/people/{uuid}", s.handleUpdatePerson()).Methods("PUT")
	s.router.HandleFunc("/people/{uuid}", s.handleDeletePerson()).Methods("DELETE")
	s.router.HandleFunc("/status", s.handleStatus()).Methods("GET")
}

// WriteHeader
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// respond ...
func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		result, _ := json.Marshal(data)
		log.Println(string(result))
		w.Write(result)
	}
}

// error response
func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// handleGetPerson bu UUID
func (s *APIServer) handleGetPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]
		person, err := s.GetPerson(uuid)
		if err != nil {
			s.error(w, r, http.StatusNotFound, errors.New("Person not found"))
			return
		}
		s.respond(w, r, http.StatusOK, person)
	}
}

// handleGetPeople
func (s *APIServer) handleGetPeople() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		people, err := s.GetPeople()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, people)
	}
}

// handleAddPerson
func (s *APIServer) handleAddPerson() http.HandlerFunc {
	type request struct {
		Survived                bool    `json:"survived"`
		PassengerClass          int     `json:"passengerClass"`
		Name                    string  `json:"name"`
		Sex                     string  `json:"sex"`
		Age                     int     `json:"age"`
		SiblingsOrSpousesAboard int     `json:"siblingsOrSpousesAboard"`
		ParentsOrChildrenAboard int     `json:"parentsOrChildrenAboard"`
		Fare                    float32 `json:"fare"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &Person{
			Survived:                req.Survived,
			PassengerClass:          req.PassengerClass,
			Name:                    req.Name,
			Sex:                     req.Sex,
			Age:                     req.Age,
			SiblingsOrSpousesAboard: req.SiblingsOrSpousesAboard,
			ParentsOrChildrenAboard: req.ParentsOrChildrenAboard,
			Fare:                    req.Fare,
		}
		if err := s.AddPerson(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, u)
	}
}

// handleUpdatePerson
func (s *APIServer) handleUpdatePerson() http.HandlerFunc {
	type request struct {
		Survived                bool    `json:"survived"`
		PassengerClass          int     `json:"passengerClass"`
		Name                    string  `json:"name"`
		Sex                     string  `json:"sex"`
		Age                     int     `json:"age"`
		SiblingsOrSpousesAboard int     `json:"siblingsOrSpousesAboard"`
		ParentsOrChildrenAboard int     `json:"parentsOrChildrenAboard"`
		Fare                    float32 `json:"fare"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &Person{
			UUID:                    uuid,
			Survived:                req.Survived,
			PassengerClass:          req.PassengerClass,
			Name:                    req.Name,
			Sex:                     req.Sex,
			Age:                     req.Age,
			SiblingsOrSpousesAboard: req.SiblingsOrSpousesAboard,
			ParentsOrChildrenAboard: req.ParentsOrChildrenAboard,
			Fare:                    req.Fare,
		}
		if err := s.UpdatePerson(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, nil)
	}
}

// handleDeletePerson
func (s *APIServer) handleDeletePerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]
		u := &Person{
			UUID: uuid,
		}
		if err := s.DeletePerson(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, nil)
	}
}

// handleStatus
func (s *APIServer) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"status": "success"})
	}
}
