package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
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
	userId := strings.TrimSpace(mux.Vars(r)["id"])

	if userId == "" {
		err := errors.New("id not given")
		s.log.Errorf("error: %v", err.Error() )
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(userId)
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

	err = model.ToJson(rw, http.StatusOK, user)
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
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = s.service.Create(context.Background(), req.Username)
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

	err = model.ToJson(rw, http.StatusOK, users)
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
	s.log.Info("Ping Request has been made....")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Pong!"))
}
