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

	success, uuid, err := h.storageClient.Register(
		ctx,
		req.Email,
		req.Password,
	)

	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Success: success,
		Uuid:    uuid,
	}, nil
}
