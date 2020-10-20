package repository

import (
	"context"
	"github.com/JamieBShaw/user-service/domain/model"
)

type Repository interface {
	UserById(ctx context.Context, id int32) (*model.User, error)
	Create(ctx context.Context, username string) error
	GetUsers(ctx context.Context)([]*model.User, error)
}
