package grpc

import (
	"context"
	"github.com/JamieBShaw/user-service/protob"
	"github.com/JamieBShaw/user-service/service"
	"github.com/sirupsen/logrus"
)


var log = logrus.New()

type grpcServer struct {
	protob.UnimplementedUserServiceServer
	service service.UserService
}

func NewGrpcServer(userService service.UserService) protob.UserServiceServer {
	return &grpcServer{
		service: userService,
	}
}

func (gs *grpcServer) GetById(ctx context.Context, req *protob.GetUserRequest) (*protob.GetUserResponse, error) {

	user, err := gs.service.GetByID(ctx, req.GetID())
	if err != nil {
		log.Errorf("error getting user by id: %v", err)
		return nil, err
	}

	res := &protob.GetUserResponse{
		User: &protob.User{
		ID:       user.ID,
		Username: user.Username,
		Admin:    user.Admin,
		},
	}

	return res, nil
}
func (gs *grpcServer) GetUsers(ctx context.Context, _ *protob.GetUsersRequest) (*protob.GetUsersResponse, error) {

	users, err := gs.service.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	req := &protob.GetUsersResponse{}

	for _, user := range users {
		u := &protob.User{
			ID:       user.ID,
			Username: user.Username,
			Admin:    user.Admin,
		}
		req.Users = append(req.Users, u)
	}

	return req, nil
}
func (gs *grpcServer) Create(ctx context.Context, req *protob.CreateUserRequest) (*protob.CreateUserResponse, error) {

	err := gs.service.Create(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &protob.CreateUserResponse{
		Confirmation: "user created",
	}, nil
}











