package http

import (
	"context"
	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/repository"
	"github.com/JamieBShaw/user-service/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	dbTest = newTestRepository()
	logTest = logrus.New()
)

func TestHandleGetUsers(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	response := httptest.NewRecorder()

	server := &httpServer{
		service: service.NewUserService(dbTest, logTest),
		router:  mux.NewRouter(),
		log:     logTest,
	}

	server.ServeHTTP(response, request)

	server.handleGetUsers(response, request)

	assertResponseBody(t, response.Body.String(), "10")

}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}




type testDB struct {}

func newTestRepository() repository.Repository {
	return &testDB{}
}

func (t testDB) UserById(ctx context.Context, id int32) (*model.User, error) {
	return nil, nil
}

func (t testDB) Create(ctx context.Context, username string) error {
	return nil
}

func (t testDB) GetUsers(ctx context.Context) ([]*model.User, error) {
	return nil, nil
}
