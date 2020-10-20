package http

import "net/http"

func (s *httpServer) routes() {
	get := s.router.Methods(http.MethodGet).Subrouter()
	post := s.router.Methods(http.MethodPost).Subrouter()

	//Get
	get.HandleFunc("/users/{id}", s.handleGetUserById)
	get.HandleFunc("/users", s.handleGetUsers)

	//Post
	post.HandleFunc("/users", s.handlePostCreateUser())
}
