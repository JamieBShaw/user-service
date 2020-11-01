package repository

import (
	"context"

	"github.com/JamieBShaw/user-service/domain/model"
)

type Repository interface {
	UserById(ctx context.Context, id int64) (*model.User, error)
	UserByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, username, password string) error
	Delete(ctx context.Context, id int64) error
	GetUsers(ctx context.Context) ([]*model.User, error)
}
