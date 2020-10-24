package service

import (
	"context"
	"errors"
	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockDb struct{}

func TestUserService_GetUserById_Test_Cases(t *testing.T) {
	tt := []struct {
		name   string
		id     int64
		res    *model.User
		errMsg string
	}{
		{
			name: "get valid user",
			id:   1,
			res: &model.User{
				ID:       1,
				Username: "David",
				Admin:    false,
			},
			errMsg: "",
		},
		{
			name:   "user does not exist",
			id:     42,
			res:    nil,
			errMsg: "could not find user with id",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := userService{
				db:  mockDb{},
				log: l,
			}
			user, err := service.GetByID(context.Background(), tc.id)
			if err != nil {
				assert.Equal(t, tc.errMsg, err.Error())
				assert.Equal(t, tc.res, user)
				return
			}
			assert.Equal(t, tc.res, user)
		})
	}
}

func TestUserService_Create_Test_Cases(t *testing.T) {
	tt := []struct {
		name     string
		username string
		errMsg   string
	}{
		{
			name:     "create user successfully",
			username: "david",
			errMsg:   "",
		},
		{
			name:     "error creating user, invalid username request",
			username: "",
			errMsg:   "username invalid",
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := userService{
				db:  mockDb{},
				log: l,
			}
			err := service.Create(context.Background(), tc.username)
			if err != nil {
				assert.Equal(t, tc.errMsg, err.Error())
				return
			}
		})
	}
}

func TestUserService_Delete_Test_Cases(t *testing.T) {
	tt := []struct {
		name   string
		id     int64
		errMsg string
	}{
		{
			name:   "User Successfully deleted",
			id:     1,
			errMsg: "",
		},
		{
			name: "Invalid request; Not valid ID",
			id: 0,
			errMsg: "invalid id",
		},
		{
			name:"Invalid request; User does not exist",
			id: 42,
			errMsg: "user not found with id",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := userService{
				db:  mockDb{},
				log: l,
			}
			err := service.Delete(context.Background(), tc.id)
			if err != nil {
				t.Logf("error: %v", err.Error())
				assert.Equal(t, tc.errMsg, err.Error())
				return
			}
		})
	}
}

func TestUserService_GetUsers(t *testing.T) {
	service := userService{
		db:  mockDb{},
		log: l,
	}
	users, err := service.GetUsers(context.Background())
	if err != nil {
		return
	}
	assert.IsType(t, []*model.User{}, users)
}

func (m mockDb) UserById(ctx context.Context, id int64) (*model.User, error) {
	users := generateUsers()

	for _, user := range users {
		if id == user.ID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m mockDb) Create(ctx context.Context, username string) error {
	if username == "" {
		return errors.New("invalid username")
	}
	return nil
}

func (m mockDb) GetUsers(ctx context.Context) ([]*model.User, error) {
	users := generateUsers()
	return users, nil
}

func (m mockDb) Delete(ctx context.Context, id int64) error {
	users := generateUsers()

	for _, user := range users {
		if id == user.ID {
			return nil
		}
	}
	return errors.New("no user found")
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
