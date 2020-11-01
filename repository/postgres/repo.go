package postgres

import (
	"context"

	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/go-pg/pg/v10"

	"github.com/sirupsen/logrus"
)

type repository struct {
	db  *pg.DB
	log *logrus.Logger
}


func NewRepository(log *logrus.Logger, db *pg.DB) *repository {
	return &repository{
		db:  db,
		log: log,
	}
}

func (repo *repository) UserById(_ context.Context, id int64) (*model.User, error) {
	repo.log.Info("[POSTGRES REPO]: Executing User By ID")

	var user model.User

	err := repo.db.Model(&user).Where("id = ?", id).First()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *repository) Create(_ context.Context, username, password string) error {
	repo.log.Info("[POSTGRES REPO]: Executing Create User")

	user := &model.User{
		Username: username,
		Password: password,
	}

	_, err := repo.db.Model(user).Insert()
	if err != nil {
		return err
	}

	if err = user.Validate(); err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetUsers(_ context.Context) ([]*model.User, error) {
	repo.log.Info("[POSTGRES REPO]: Executing Get Users")

	var users []*model.User

	err := repo.db.Model(&users).Select()
	if err != nil {
		repo.log.Errorf("error selecting users: %v", err)
		return nil, err
	}

	return users, nil
}

func (repo *repository) Delete(ctx context.Context, id int64) error {
	repo.log.Info("[POSTGRES REPO]: Executing Delete User")

	user := &model.User{
		ID: id,
	}
	_, err := repo.db.Model(user).Where("id = ?", id).Delete()
	if err != nil {
		repo.log.Errorf("error deleting user: %v", err)
		return err
	}

	return nil
}

func (repo *repository) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	repo.log.Info("[POSTGRES REPO]: Executing Getting User by Username")

	user := &model.User{
		Username: username,
	}

	err := repo.db.Model(user).Where("username = ?", username).First()
	if err != nil {
		repo.log.Errorf("error getting user by username: %s, error: %v",username, err )
		return nil, err
	}

	return user, nil
}

