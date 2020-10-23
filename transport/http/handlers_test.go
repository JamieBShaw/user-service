package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/repository"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	db repository.Repository
}

type mockDb struct{}

func TestHttpServer_GetUsers_Valid_Response(t *testing.T) {
	serverMock := httpServer{
		service: mockUserService{},
	}
	req, err := http.NewRequest("GET", "localhost:50051/users", nil)
	if err != nil {
		t.Fatalf("could not create mock request: %v", err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(serverMock.GetUsers)

	handler.ServeHTTP(rec, req)
	res := rec.Result()
	users := marshalResponse(res)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.IsType(t, []*model.User{}, users)
}

func TestHttpServer_GetById_Test_Cases(t *testing.T) {
	tt := []struct {
		name   string
		status int
		expectedRes string
		errMsg string
		userId string
	}{
		{
			name:   "valid user response",
			status: 200,
			errMsg: "",
			expectedRes: "{\"id\":1,\"username\":\"David\",\"admin\":false}",
			userId: "1",
		},
		{
			name:   "user with that id does not exist",
			status: 404,
			expectedRes: "",
			errMsg: "could not find user with id",
			userId: "42",
		},
		{
			name: "invalid url parameter",
			status: 500,
			expectedRes: "",
			errMsg: "invalid query parameter",
			userId: "foxtrot",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serverMock := httpServer{
				service: mockUserService{},
				log: logrus.New(),
			}
			req, err := http.NewRequest("GET", "localhost:50051/users/?id=" + tc.userId, nil)
			if err != nil {
				t.Fatalf("could not create mock request: %v", err)
			}
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(serverMock.GetById)

			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.errMsg == "" {
				// Happy Path
				assert.Equal(t, tc.status, res.StatusCode)
				assert.Equal(t, tc.expectedRes, string(bytes.TrimSpace(b)))
				return
			}
			assert.Equal(t, tc.status, res.StatusCode)
			assert.Equal(t, tc.errMsg, string(bytes.TrimSpace(b)))

		})

	}

}

func TestHttpServer_Ping(t *testing.T) {
	server := httpServer{}
	req, err := http.NewRequest("GET", "localhost:50051/ping", nil)
	if err != nil {
		t.Fatalf("could not create mock request: %v", err)
	}
	rec := httptest.NewRecorder()

	server.Ping(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read respinse: %v", err)
	}

	assert.Equal(t, "Pong!", string(b))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func (m mockUserService) GetByID(_ context.Context, id int64) (*model.User, error) {
	users := generateUsers()

	for _, user := range users {
		if id == user.ID {
			return user, nil
		}
	}
	return nil, errors.New("could not find user with id")
}

func (m mockUserService) GetUsers(ctx context.Context) ([]*model.User, error) {
	users := generateUsers()
	return users, nil
}

func (m mockUserService) Create(ctx context.Context, username string) error {
	return nil
}

func generateUsers() []*model.User {
	var users []*model.User
	names := []string{"James", "David", "Michael"}

	for i, name := range names {
		users = append(users, &model.User{
			ID:        int64(i),
			Username:  name,
			Admin:     false,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		})
	}

	return users
}

func marshalResponse(res *http.Response) []*model.User {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	var users []*model.User

	err = json.Unmarshal(b, &users)
	if err != nil {
		return nil
	}

	return users
}
