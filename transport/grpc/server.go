package grpc

import (
	"context"
	"github.com/JamieBShaw/user-service/protob"
	"github.com/JamieBShaw/user-service/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if req == nil || req.GetID() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument")
	}

	user, err := gs.service.GetByID(ctx, req.GetID())
	if err != nil {
		log.Errorf("error getting user by id: %v", req.GetID())
		return nil, status.Errorf(codes.NotFound, "user not found")
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

func (gs *grpcServer) GetUsers(ctx context.Context, req *protob.GetUsersRequest) (*protob.GetUsersResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	users, err := gs.service.GetUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "users not found")
	}

	res := &protob.GetUsersResponse{}

	for _, user := range users {
		u := &protob.User{
			ID:       user.ID,
			Username: user.Username,
			Admin:    user.Admin,
		}
		res.Users = append(res.Users, u)
	}

	return res, nil
}

func (gs *grpcServer) Create(ctx context.Context, req *protob.CreateUserRequest) (*protob.CreateUserResponse, error) {
	if req == nil || req.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	err := gs.service.Create(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &protob.CreateUserResponse{
		Confirmation: "user created",
	}, nil
}

func (gs *grpcServer) Delete(ctx context.Context, req *protob.DeleteUserRequest) (*protob.DeleteUserResponse, error) {
	if req == nil || req.GetID() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	err := gs.service.Delete(ctx, req.GetID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &protob.DeleteUserResponse{
		Confirmation: "user deleted",
	}, nil
}
