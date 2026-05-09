package handler

import (
	"context"
	"log"
	"storage-service/internal/config"

	pb "storagepb"
)

type StorageHandler struct {
	pb.UnimplementedStorageServiceServer
	envConf *config.Config
}

func NewStorageHandler(envConf *config.Config) *StorageHandler {
	return &StorageHandler{envConf: envConf}
}

func (h *StorageHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register called: email=%s", req.Email)

	return &pb.RegisterResponse{
		Success: true,
		Uuid:    "good",
	}, nil
}
