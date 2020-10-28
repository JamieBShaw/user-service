package http

import "net/http"

func (s *httpServer) routes() {
	get := s.router.Methods(http.MethodGet).Subrouter()
	post := s.router.Methods(http.MethodPost).Subrouter()
	deleteR := s.router.Methods(http.MethodDelete).Subrouter()
	//Get
	get.HandleFunc("/users/{id}", s.GetById)
	get.HandleFunc("/users", s.GetUsers)
	//Post
	post.HandleFunc("/users", s.Create())
	//Delete
	deleteR.HandleFunc("/users/{id}", s.Delete)
	//PING
	get.HandleFunc("/healthz", s.Healthz)
}
