package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/JamieBShaw/user-service/protob"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/repository"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	db repository.Repository
}

type mockAuthClient struct{}

func TestHttpServer_GetUsers_Valid_Response(t *testing.T) {
	serverMock := httpServer{
		service: mockUserService{},
		log:     l,
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
		name        string
		status      int
		expectedRes string
		errMsg      string
		userId      string
	}{
		{
			name:        "valid user response",
			status:      200,
			errMsg:      "",
			expectedRes: "{\"id\":1,\"username\":\"James\",\"password\":\"password\",\"admin\":false}",
			userId:      "1",
		},
		{
			name:        "user with that id does not exist",
			status:      404,
			expectedRes: "",
			errMsg:      "could not find user with id",
			userId:      "42",
		},
		{
			name:        "invalid url parameter, invalid id given",
			status:      500,
			expectedRes: "",
			errMsg:      "invalid query parameter",
			userId:      "foxtrot",
		},
		{
			name:        "invalid url parameter, no id given",
			status:      500,
			expectedRes: "",
			errMsg:      "id not given",
			userId:      "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			serverMock := httpServer{
				service: mockUserService{},
				log:     l,
			}
			req, err := http.NewRequest("GET", "localhost:50051/users/", nil)
			if err != nil {
				t.Fatalf("could not create mock request: %v", err)
			}
			req = mux.SetURLVars(req, map[string]string{
				"id": tc.userId,
			})
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

func TestHttpServer_Create(t *testing.T) {
	tt := []struct {
		name        string
		username    string
		status      int
		expectedRes string
		errMsg      string
	}{
		{
			name:        "valid user created",
			username:    "David",
			status:      201,
			expectedRes: "User successfully created",
			errMsg:      "",
		},
		{
			name:        "invalid request",
			username:    "1234",
			status:      500,
			expectedRes: "",
			errMsg:      "user not created",
		},
		{
			name:        "invalid request, username too long",
			username:    "longusernameover10",
			status:      500,
			expectedRes: "",
			errMsg:      "user not created",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			serverMock := httpServer{
				service: mockUserService{},
				log:     l,
			}
			body := strings.NewReader("{\n\"username\": \"" + tc.username + "\"\n}")

			req, err := http.NewRequest("POST", "localhost:50051/users/", body)
			if err != nil {
				t.Fatalf("could not create mock request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := serverMock.Register()

			handler.ServeHTTP(rec, req)

			res := rec.Result()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.errMsg == "" {
				// Good Path
				assert.Equal(t, tc.expectedRes, string(b))
				assert.Equal(t, tc.status, res.StatusCode)
				return
			}
			assert.Equal(t, tc.errMsg, string(bytes.TrimSpace(b)))
			assert.Equal(t, tc.status, res.StatusCode)
		})
	}
}

func TestHttpServer_Delete_Test_Cases(t *testing.T) {
	tt := []struct {
		name   string
		id     string
		res    string
		errMsg string
		status int
	}{
		{
			name:   "User successfully deleted",
			id:     "1",
			res:    "User successfully deleted",
			errMsg: "",
			status: 200,
		},
		{
			name:   "Invalid request, id not convertible to int",
			id:     "notAnInt",
			res:    "",
			errMsg: "invalid query parameter",
			status: 500,
		},
		{
			name:   "Invalid request, id not allowed",
			id:     "0",
			res:    "",
			errMsg: "invalid id",
			status: 404,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			serverMock := httpServer{
				service: mockUserService{},
				log:     l,
			}

			req, err := http.NewRequest("POST", "localhost:50051/users/", nil)
			if err != nil {
				t.Fatalf("could not create mock request: %v", err)
			}
			req = mux.SetURLVars(req, map[string]string{
				"id": tc.id,
			})

			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(serverMock.Delete)
			handler.ServeHTTP(rec, req)

			res := rec.Result()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.errMsg == "" {
				// Good Path
				assert.Equal(t, tc.res, string(bytes.TrimSpace(b)))
				assert.Equal(t, tc.status, res.StatusCode)
				return
			}
			assert.Equal(t, tc.errMsg, string(bytes.TrimSpace(b)))
			assert.Equal(t, tc.status, res.StatusCode)
		})
	}
}

func TestHttpServer_Healthz(t *testing.T) {
	server := httpServer{
		nil,
		nil,
		l,
		nil,
	}
	req, err := http.NewRequest("GET", "localhost:50051/ping", nil)
	if err != nil {
		t.Fatalf("could not create mock request: %v", err)
	}
	rec := httptest.NewRecorder()

	server.Healthz(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read respinse: %v", err)
	}

	assert.Equal(t, "Healthy!", string(b))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

//
//func TestHttpServer_Login(t *testing.T) {
//	tt := []struct {
//		name string
//		username string
//		password string
//		errMsg string
//		status int
//		response string
//	}{
//		{
//			name: "User successfuly recieved",
//			username: "James",
//			password: "password",
//			errMsg: "",
//			status: 200,
//			response: "{\"access_token\":\"3214343254\",\"refresh_token\":\"5435436265\"}",
//		},
//		{
//			name: "User not found",
//			username: "",
//			password: "password",
//			errMsg: "user not found with username",
//			status: 500,
//			response: "",
//		},
//		{
//			name: "invalid username",
//			username: "ia1",
//			password: "password",
//			errMsg: "invalid username",
//			status: 500,
//			response: "",
//		},
//	}
//	for _, tc := range tt {
//		tc := tc
//
//		type UserLoginRequest struct {
//			Username string `json:"username"`
//			Password string `json:"password"`
//		}
//
//		t.Run(tc.name, func(t *testing.T) {
//			t.Parallel()
//			serverMock := httpServer{
//				service: mockUserService{},
//				log:     l,
//				authServiceClient: mockAuthClient{},
//			}
//
//			loginReq := &UserLoginRequest{
//				Username: tc.username,
//				Password: tc.password,
//			}
//
//			reqByte, _ := json.Marshal(loginReq)
//
//			body := bytes.NewReader(reqByte)
//
//			req, err := http.NewRequest("POST", "localhost:50051/login", body)
//			if err != nil {
//				t.Fatalf("could not create mock request: %v", err)
//			}
//
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			handler := serverMock.Login()
//			handler.ServeHTTP(rec, req)
//
//			res := rec.Result()
//
//			b, err := ioutil.ReadAll(res.Body)
//			if err != nil {
//				t.Fatalf("could not read response: %v", err)
//			}
//
//			if tc.errMsg == "" {
//				// Good Path
//				assert.Equal(t, tc.response, string(bytes.TrimSpace(b)))
//				assert.Equal(t, tc.status, res.StatusCode)
//				return
//			}
//			assert.Equal(t, tc.errMsg, string(b))
//			assert.Equal(t, tc.status, res.StatusCode)
//		})
//	}
//
//}

func (m mockUserService) GetByID(_ context.Context, id int64) (*model.User, error) {
	users := generateUsers()

	for _, user := range users {
		if id == user.ID {
			return user, nil
		}
	}
	return nil, errors.New("could not find user with id")
}

func (m mockUserService) GetUsers(_ context.Context) ([]*model.User, error) {
	users := generateUsers()
	return users, nil
}

func (m mockUserService) Create(_ context.Context, username, password string) error {
	users := generateUsers()

	for _, user := range users {
		if username == user.Username {
			return nil
		}
	}
	return errors.New("user not created")
}

func (m mockUserService) Delete(_ context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	users := generateUsers()
	for _, user := range users {
		if id == user.ID {
			return nil
		}
	}
	return errors.New("user does not exist")
}

func (m mockUserService) GetByUsernameAndPassword(_ context.Context, username, password string) (*model.User, error) {
	users := generateUsers()

	for _, user := range users {
		if username == user.Username {
			log.Printf("username: %v", username)
			return user, nil
		}
	}

	return nil, errors.New("error username")
}

func (m mockAuthClient) CreateAccessToken(ctx context.Context, in *protob.CreateAccessTokenRequest, opts ...grpc.CallOption) (*protob.CreateAccessTokenResponse, error) {
	return &protob.CreateAccessTokenResponse{
		AuthToken:    "3214343254",
		RefreshToken: "5435436265",
	}, nil
}

func (m mockAuthClient) DeleteAccessToken(ctx context.Context, in *protob.DeleteAccessTokenRequest, opts ...grpc.CallOption) (*protob.DeleteAccessTokenResponse, error) {
	return &protob.DeleteAccessTokenResponse{Confirmation: "Success"}, nil
}

func generateUsers() []*model.User {
	var users []*model.User
	names := []string{"James", "David", "Michael", "jimmy", "michael", "teddy", "maclom"}
	passwords := []string{"password", "password123", "laskdkad", "djpsafd", "jdfpsajf", "dpsa111", "dlwfops"}

	for i, name := range names {
		for _, password := range passwords {
			users = append(users, &model.User{
				ID:        int64(i + 1),
				Username:  name,
				Admin:     false,
				Password:  password,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			})
		}
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
