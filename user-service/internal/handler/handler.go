package handler

import (
	"context"
	"log"
	"user-service/internal/config"
	grpcclient "user-service/internal/grpc-client"

	pb "userpb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	envConf       *config.Config
	storageClient *grpcclient.StorageClient
}

func NewUserHandler(envConf *config.Config, storageClient *grpcclient.StorageClient) *UserHandler {
	return &UserHandler{envConf: envConf, storageClient: storageClient}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register called: email=%s", req.Email)

	return &pb.RegisterResponse{
		Success: true,
		Message: "Registration successful (stub)",
	}, nil
}
