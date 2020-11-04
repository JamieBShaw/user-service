package auth_serivce_grpc

import (
	googlegrpc "google.golang.org/grpc"
	"log"
)


func NewAuthClientConn() (* googlegrpc.ClientConn) {
	cc, err := googlegrpc.Dial("0.0.0.0:8081", googlegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("error getting connection grpc client: %v", err)
		return nil
	}

	log.Println("Starting GRPC AuthService Client on port: 8081 ..........")
	return cc
}
