package auth_serivce_grpc

import (

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

func NewAuthClientConn() *grpc.ClientConn {
	cc, err := grpc.Dial("0.0.0.0:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error getting connection grpc client: %v", err)
		return nil
	}

	log.Info("Starting GRPC AuthService Client on port: 8081.")
	return cc
}
