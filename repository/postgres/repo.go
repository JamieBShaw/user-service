package postgres

import (
	"context"
	"github.com/JamieBShaw/user-service/domain/model"
	"github.com/go-pg/pg/v10"

	"github.com/sirupsen/logrus"
)

type repository struct {
	db *pg.DB
	log *logrus.Logger
}

func NewRepository(log *logrus.Logger, db *pg.DB) *repository {
		return &repository{
			db: db,
			log: log,
		}
}

func (repo *repository) UserById(_ context.Context, id int64) (*model.User, error) {
	repo.log.Info("REPO: Executing User By ID")

	var user model.User

	err := repo.db.Model(&user).Where("id = ?", id).First()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *repository) Create(_ context.Context, username string) error {
	repo.log.Info("REPO: Executing Create User")

	user := &model.User{
		Username: username,
	}

	_, err := repo.db.Model(&user).Returning("*").Insert()
	if err != nil {
		return err
	}

	if err = user.Validate(); err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetUsers(_ context.Context)([]*model.User, error){
	repo.log.Info("REPO: Executing Get Users")

	var users []*model.User

	err := repo.db.Model(&users).Select()
	if err != nil {
		repo.log.Errorf("error selecting users: %v", err)
		return nil, err
	}

	return users, nil
}


