package http

import (
	"net/http"

	"github.com/JamieBShaw/user-service/protob"
	"github.com/JamieBShaw/user-service/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	l = logrus.New()
)

type Server interface {
	GetById(rw http.ResponseWriter, r *http.Request)
	GetUsers(rw http.ResponseWriter, r *http.Request)
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
	Delete(rw http.ResponseWriter, r *http.Request)
	Healthz(rw http.ResponseWriter, r *http.Request)
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

type httpServer struct {
	service           service.UserService
	router            *mux.Router
	log               *logrus.Logger
	authServiceClient protob.AuthServiceClient
}

func (s *httpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func NewHttpHandler(service service.UserService, router *mux.Router, client protob.AuthServiceClient) http.Handler {
	server := &httpServer{service, router, l, client}
	server.routes()

	return server
}
