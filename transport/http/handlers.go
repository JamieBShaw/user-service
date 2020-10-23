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
	Ping(rw http.ResponseWriter, r *http.Request)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

func (s *httpServer) GetById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, errors.New("invalid query parameter").Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	user, err := s.service.GetByID(ctx, int64(id))
	if err != nil {
		s.log.Errorf("error: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(rw).Encode(user)
}

func (s *httpServer) Create() http.HandlerFunc {

	type request struct {
		Username string
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		ctx := context.Background()
		err = s.service.Create(ctx, req.Username)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte("User successfully created"))
	}
}

func (s *httpServer) GetUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := context.Background()

	users, err := s.service.GetUsers(ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

	}
}

func (s *httpServer) Ping(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Pong!"))
}
