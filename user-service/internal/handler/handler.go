package handler

import (
	"context"
	"log"
	"user-service/internal/config"
	"user-service/internal/service"

	pb "userpb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	envConf *config.Config
	service *service.Service
}

func NewUserHandler(envConf *config.Config, service *service.Service) *UserHandler {
	return &UserHandler{envConf: envConf, service: service}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register called: email=%s", req.Email)

	exists, err := h.service.User.IsEmailExists(req.Email)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Uuid:    "",
		}, nil
	}
	if exists {
		log.Printf("Email already registered: %s", req.Email)
		return &pb.RegisterResponse{
			Success: false,
			Uuid:    "",
		}, nil
	}

	user, err := h.service.User.Add(req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Uuid:    "",
		}, nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Uuid:    user.Uuid,
	}, nil
}
