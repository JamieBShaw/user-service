package service

import (
	"context"
	"errors"

	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/repository"
	"github.com/sirupsen/logrus"
)

var (
	l = logrus.New()
)

type userService struct {
	db  repository.Repository
	log *logrus.Logger
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsernameAndPassword(ctx context.Context, username, password string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	Create(ctx context.Context, username, password string) error
	Delete(ctx context.Context, id int64) error
}

func NewUserService(db repository.Repository) *userService {
	return &userService{db: db, log: l}
}

func (u *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u.log.Info("[USER SERVICE]: Get User by ID")

	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	user, err := u.db.UserById(ctx, id)
	if err != nil {
		u.log.Errorf("USER SERVICE: error: %v", err)
		return nil, errors.New("could not find user with id")
	}

	err = user.Validate()
	if err != nil {
		u.log.Errorf("unable to validate user: %v", user)
		return nil, errors.New("could not validate user")
	}

	return user, nil
}

func (u *userService) Create(ctx context.Context, username, password string) error {
	u.log.Info("[USER SERVICE]: Register User:" + username)

	if username == "" || len(username) > 10 {
		return errors.New("username invalid")
	}

	if password == "" || len(password) < 8 {
		return errors.New("password invalid")
	}

	err := u.db.Create(ctx, username, password)
	if err != nil {
		return errors.New("error creating user")
	}

	return nil
}

func (u *userService) GetUsers(ctx context.Context) ([]*model.User, error) {
	u.log.Info("[USER SERVICE]: Get Users")

	users, err := u.db.GetUsers(ctx)
	if err != nil {
		return nil, errors.New("unable to get users")
	}

	for _, user := range users {
		err = user.Validate()
		if err != nil {
			u.log.Errorf("error user in users array not valid: %d, error: %v", user.ID, err)
		}
	}

	return users, nil
}

func (u *userService) Delete(ctx context.Context, id int64) error {
	u.log.Info("[USER SERVICE]: Delete User")

	if id <= 0 {
		return errors.New("invalid id")
	}

	err := u.db.Delete(ctx, id)
	if err != nil {
		return errors.New("user not found with id")
	}

	return nil
}

func (u *userService) GetByUsernameAndPassword(ctx context.Context, username, password string) (*model.User, error) {
	u.log.Info("[USER SERVICE]: Get User by Username")

	user, err := u.db.UserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found with username")
	}

	err = user.ValidatePassword(password)
	if err != nil {
		return nil, errors.New("user password incorrect")
	}

	err = user.Validate()
	if err != nil {
		u.log.Errorf("unable to validate user: %v", user)
		return nil, errors.New("could not validate user")
	}

	return user, nil
}
