package handler

import (
	"context"
	"log"

	pb "userpb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register called: email=%s", req.Email)

	return &pb.RegisterResponse{
		Success: true,
		Message: "Registration successful (stub)",
	}, nil
}
