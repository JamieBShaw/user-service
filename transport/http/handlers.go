package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *httpServer) handleGetUserById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	user, err := s.service.GetByID(ctx, int32(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(rw).Encode(user)
}

func (s *httpServer) handlePostCreateUser() http.HandlerFunc {

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

func (s *httpServer) handleGetUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := context.Background()

	users, err := s.service.GetUsers(ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(users)
}
