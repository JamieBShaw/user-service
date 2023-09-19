package http

import (
	"net/http"

	auth "github.com/JamieBShaw/auth-service/middleware/http_middleware"
)

func (s *httpServer) routes() {

	get := s.router.Methods(http.MethodGet).Subrouter()
	post := s.router.Methods(http.MethodPost).Subrouter()
	deleteR := s.router.Methods(http.MethodDelete).Subrouter()
	//Get
	get.HandleFunc("/users/{id}", s.GetById)
	get.HandleFunc("/users", s.GetUsers)
	//Post
	post.HandleFunc("/register", s.Register())
	post.HandleFunc("/login", s.Login())
	post.HandleFunc("/logout", auth.AuthenticationMiddleware(s.Logout()))
	//Delete
	deleteR.HandleFunc("/users/{id}", s.Delete)
	//PING
	get.HandleFunc("/healthz", s.Healthz)

}
