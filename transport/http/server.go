package http

import (
	"github.com/JamieBShaw/user-service/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type httpServer struct {
	service service.UserService
	router *mux.Router
	log *logrus.Logger
}

func (s *httpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func NewHttpHandler(service service.UserService, router *mux.Router, log *logrus.Logger) http.Handler {
	server := &httpServer{service, router, log}
	server.routes()

	return server
}