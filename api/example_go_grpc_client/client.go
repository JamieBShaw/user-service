package main

import (
	"context"
	"fmt"
	"github.com/JamieBShaw/user-service/protob"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// Example grpc client I built before attempting to implement
func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error getting connection grpc example_go_grpc_client: %v", err)
	}

	log.Println("listening on port 50051")

	defer cc.Close()

	c := protob.NewUserServiceClient(cc)

	getUsers(c)

}

func getUserByID(c protob.UserServiceClient, req *protob.GetUserRequest) {

	res, err := c.GetById(context.Background(), req)
	if err != nil {
		err, ok := status.FromError(err)
		if ok {
			log.Printf("Message: %v , Code: %v", err.Message(), err.Code())

			if err.Code() == codes.InvalidArgument {
				log.Printf("user name not valid: %v", err)
			}

		} else {
			log.Fatalf("error not grpc error: %v", err)
		}
	}
	fmt.Printf("USER: %v", res.GetUser())
}

func getUsers(c protob.UserServiceClient) {

	req := &protob.GetUsersRequest{}

	res, err := c.GetUsers(context.Background(), req)
	if err != nil {
		err, ok := status.FromError(err)
		if ok {
			log.Printf("Message: %v, Code: %v", err.Message(), err.Code())

			if err.Code() == codes.InvalidArgument {
				log.Printf("user name not valid: %v", err)
			}
		} else {
			log.Fatalf("error not grpc error: %v", err)
		}
	}

	users := res.GetUsers()

	for _, user := range users {
		log.Printf("User: %v", user)
	}

}

func createUser(c protob.UserServiceClient) {

	req := &protob.CreateUserRequest{Username: "Jamie0001"}

	res, err := c.Create(context.Background(), req)
	if err != nil {
		err, ok := status.FromError(err)
		if ok {
			log.Printf("Message: %v, Code: %v", err.Message(), err.Code())

			if err.Code() == codes.InvalidArgument {
				log.Printf("user name not valid: %v", err)
			}
		} else {
			log.Fatalf("error not grpc error: %v", err)
		}
	}

	log.Printf("response: %v", res.GetConfirmation())

}
