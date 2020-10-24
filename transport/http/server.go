package http

import (
	"github.com/JamieBShaw/user-service/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	l = logrus.New()
)

type httpServer struct {
	service service.UserService
	router  *mux.Router
	log     *logrus.Logger
}

func (s *httpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func NewHttpHandler(service service.UserService, router *mux.Router) http.Handler {
	server := &httpServer{service, router, l}
	server.routes()

	return server
}
