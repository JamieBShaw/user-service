package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/protob"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

type mockUserService struct {
}

func TestGrpcServer_GetById_Test_Cases(t *testing.T) {
	tt := []struct {
		name     string
		id       int64
		errMsg   string
		response *protob.User
		errCode  string
	}{
		{
			name:   "get valid user",
			id:     1,
			errMsg: "",
			response: &protob.User{
				ID:       1,
				Username: "David",
				Admin:    false,
			},
		},
		{
			name:     "user not found with id",
			id:       42,
			errMsg:   "user not found",
			response: nil,
			errCode:  "NotFound",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := grpcServer{service: mockUserService{}}
			res, err := server.GetById(context.Background(), &protob.GetUserRequest{ID: tc.id})
			if err != nil {
				statusErr, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tc.errCode, statusErr.Code().String())
					assert.Equal(t, tc.errMsg, statusErr.Message())
					assert.Equal(t, tc.response, res.GetUser())
				}
				return
			}
			assert.Equal(t, tc.response, res.GetUser())
		})
	}
}

func TestGrpcServer_GetUsers(t *testing.T) {
	tt := []struct {
		name     string
		req      *protob.GetUsersRequest
		errMsg   string
		response []*protob.User
		errCode  string
	}{
		{
			name:     "valid request, users returned",
			req:      &protob.GetUsersRequest{},
			errMsg:   "",
			response: generateProtoUsers(),
		},
		{
			name:     "invalid request, no users returned",
			req:      nil,
			errMsg:   "invalid request",
			response: nil,
			errCode:  "InvalidArgument",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := grpcServer{service: mockUserService{}}
			res, err := server.GetUsers(context.Background(), tc.req)
			if err != nil {
				statusErr, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tc.errCode, statusErr.Code().String())
					assert.Equal(t, tc.errMsg, statusErr.Message())
					assert.Equal(t, tc.response, res.GetUsers())
				}
				return
			}
			assert.Equal(t, tc.response, res.GetUsers())
		})
	}
}

func TestGrpcServer_Create_Test_Cases(t *testing.T) {
	tt := []struct {
		name   string
		req    *protob.CreateUserRequest
		res    string
		errMsg string
		code   string
	}{
		{
			name: "request good, create valid user",
			req: &protob.CreateUserRequest{
				Username: "James",
			},
			res:    "user created",
			errMsg: "",
			code:   "",
		},
		{
			name:   "bad request, username invalid",
			req:    nil,
			res:    "",
			errMsg: "invalid request",
			code:   "InvalidArgument",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := grpcServer{service: mockUserService{}}
			res, err := server.Create(context.Background(), tc.req)
			if err != nil {
				statusErr, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tc.code, statusErr.Code().String())
					assert.Equal(t, tc.errMsg, statusErr.Message())
					assert.Equal(t, tc.res, res.GetConfirmation())
				}
				return
			}
			assert.Equal(t, tc.res, res.GetConfirmation())
		})
	}
}

func TestGrpcServer_Delete_Test_Cases(t *testing.T) {
	tt := []struct {
		name   string
		req    *protob.DeleteUserRequest
		res    string
		errMsg string
		code   string
	}{
		{
			name: "User deleted successfully",
			req: &protob.DeleteUserRequest{
				ID: 1,
			},
			res:    "user deleted",
			errMsg: "",
		},
		{
			name: "User not found with that id",
			req: &protob.DeleteUserRequest{
				ID: 42,
			},
			res:    "",
			errMsg: "user not found",
			code:   "Internal",
		},
		{
			name:   "invalid request, nil request",
			req:    nil,
			res:    "",
			errMsg: "invalid request",
			code:   "InvalidArgument",
		},
		{
			name:   "invalid request, invalid id",
			req:    &protob.DeleteUserRequest{ID: 0},
			res:    "",
			errMsg: "invalid request",
			code:   "InvalidArgument",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := grpcServer{service: mockUserService{}}
			res, err := server.Delete(context.Background(), tc.req)
			if err != nil {
				statusErr, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tc.code, statusErr.Code().String())
					assert.Equal(t, tc.errMsg, statusErr.Message())
					assert.Equal(t, tc.res, res.GetConfirmation())
				}
			}

			assert.Equal(t, tc.res, res.GetConfirmation())
		})
	}

}

func (m mockUserService) GetByUsernameAndPassword(ctx context.Context, username, password string) (*model.User, error) {
	panic("implement me")
}

func (m mockUserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	users := generateUsers()
	for _, user := range users {
		if id == user.ID {
			return user, nil
		}
	}
	return nil, errors.New("error: getting user by Id")
}

func (m mockUserService) GetUsers(ctx context.Context) ([]*model.User, error) {
	return generateUsers(), nil
}

func (m mockUserService) Create(ctx context.Context, username, password string) error {
	if username == "James" {
		return nil
	}
	return errors.New("cannot create user")
}
func (m mockUserService) Delete(ctx context.Context, id int64) error {
	users := generateUsers()

	for _, user := range users {
		if id == user.ID {
			return nil
		}
	}
	return errors.New("user not found")
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

func generateProtoUsers() []*protob.User {
	var users []*protob.User
	names := []string{"James", "David", "Michael"}

	for i, name := range names {
		users = append(users, &protob.User{
			ID:       int64(i),
			Username: name,
			Admin:    false,
		})
	}
	return users
}
