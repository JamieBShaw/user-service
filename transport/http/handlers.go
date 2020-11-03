package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/protob"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const CurrentUserKey = "currentUser"



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

func (s *httpServer) Register() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		s.log.Info("[HTTP SERVER]: Executing Register Handler")
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

func (s *httpServer) Login() http.HandlerFunc {

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

		user, err := s.service.GetByUsernameAndPassword(context.Background(), req.Username, req.Password)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		res , err := s.client.CreateAccessToken(ctx, &protob.CreateAccessTokenRequest{
			ID: user.ID,
		})
		if err != nil {
			s.log.Errorf("error: %v", err)

			err, ok := status.FromError(err)
			if ok {
				if err.Code() == codes.DeadlineExceeded {
					s.log.Errorf("timeout exceeded: %v", err)
					return
				}
				s.log.Errorf("unexpected error: %v", err)
			}
			s.log.Errorf("error whilke calling CreateAccessToken RPC: %v", err)
			return
		}
		tokens := map[string]string{
			"access_token":  res.GetAuthToken(),
			"refresh_token": res.GetRefreshToken(),
		}

		err = json.NewEncoder(rw).Encode(&tokens)
		if err != nil {
			s.log.Errorf("error: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func (s *httpServer) Logout() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		s.log.Info("[HTTP SERVER]: Executing Login Handler")
	}
}


func (s *httpServer) Healthz(rw http.ResponseWriter, r *http.Request) {
	s.log.Info("Ping Request has been made....")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Healthy!"))
}



