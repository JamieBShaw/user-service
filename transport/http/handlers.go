package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Server interface {
	GetById(rw http.ResponseWriter, r *http.Request)
	GetUsers(rw http.ResponseWriter, r *http.Request)
	Create() http.HandlerFunc
	Delete(rw http.ResponseWriter, r *http.Request)
	Ping(rw http.ResponseWriter, r *http.Request)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

func (s *httpServer) GetById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, errors.New("invalid query parameter").Error(), http.StatusInternalServerError)
		return
	}

	user, err := s.service.GetByID(context.Background(), int64(id))
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusNotFound)
	}
}

func (s *httpServer) Create() http.HandlerFunc {

	type request struct {
		Username string `json:"username"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		s.log.Infof("request: %v", req.Username)

		ctx := context.Background()
		err = s.service.Create(ctx, req.Username)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("User successfully created"))
	}
}

func (s *httpServer) GetUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := context.Background()

	users, err := s.service.GetUsers(ctx)
	if err != nil {
		s.log.Errorf("error: %v", err)
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(rw).Encode(users)
	if err != nil {
		s.log.Errorf("error: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (s *httpServer) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, errors.New("invalid query parameter").Error(), http.StatusInternalServerError)
		return
	}

	err = s.service.Delete(context.Background(), int64(id))
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User successfully deleted"))
}

func (s *httpServer) Ping(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Pong!"))
}
