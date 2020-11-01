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

const CurrentUserKey = "currentUser"

type Server interface {
	GetById(rw http.ResponseWriter, r *http.Request)
	GetUsers(rw http.ResponseWriter, r *http.Request)
	Create() http.HandlerFunc
	Login() http.HandlerFunc
	Delete(rw http.ResponseWriter, r *http.Request)
	Healthz(rw http.ResponseWriter, r *http.Request)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

func (s *httpServer) GetById(rw http.ResponseWriter, r *http.Request) {
	s.log.Info("[HTTP SERVER]: Executing GetById Handler")
	userId := strings.TrimSpace(mux.Vars(r)["id"])

	if userId == "" {
		err := errors.New("id not given")
		s.log.Errorf("error: %v", err.Error())
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
		Password string `json:"password"`
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		s.log.Info("[HTTP SERVER]: Executing Create Handler")
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = s.service.Create(context.Background(), req.Username, req.Password)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("User successfully created"))
	}
}

func (s *httpServer) GetUsers(rw http.ResponseWriter, r *http.Request) {
	s.log.Info("[HTTP SERVER]: Executing GetUsers Handler")
	users, err := s.service.GetUsers(context.Background())
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
	s.log.Info("[HTTP SERVER]: Executing Delete Handler")
	userId := strings.TrimSpace(mux.Vars(r)["id"])

	id, err := strconv.Atoi(userId)
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

func (s *httpServer) Login(next http.HandlerFunc) http.HandlerFunc {

	type UserLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		s.log.Info("[HTTP SERVER]: Executing Login Handler")
		var req UserLoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		user, err := s.service.GetByUsername(context.Background(), req.Username)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = user.Validate()
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if user.Password != req.Password {
			http.Error(rw, errors.New("invalid credentials").Error(), http.StatusForbidden)
		}

		ctx := context.WithValue(r.Context(), CurrentUserKey, user.ID )
		r = r.WithContext(ctx)
		next(rw, r)
	}
}

func (s *httpServer) Healthz(rw http.ResponseWriter, r *http.Request) {
	s.log.Info("Ping Request has been made....")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Pong!"))
}



