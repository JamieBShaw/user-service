package service

import (
	"context"
	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/JamieBShaw/user-service/repository"
	"github.com/sirupsen/logrus"
)


type userService struct {
	db repository.Repository
	log *logrus.Logger
}

type UserService interface {
	GetByID(ctx context.Context, id int32) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	Create(ctx context.Context, username string) error
}

func NewUserService(db repository.Repository, log *logrus.Logger) *userService {
	return &userService{db: db, log: log}
}

func (u *userService) GetByID(ctx context.Context, id int32) (*model.User, error) {
	u.log.Info("USER SERVICE: Get User by ID")

	user, err := u.db.UserById(ctx, id)
	if err != nil {
		u.log.Errorf("USER SERVICE: error: %v", err)
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) Create(ctx context.Context, username string) error {
	u.log.Info("USER SERVICE: Create User")

	err := u.db.Create(ctx, username)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) GetUsers(ctx context.Context) ([]*model.User, error) {
	u.log.Info("USER SERVICE: Get Users")

	users, err := u.db.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}




