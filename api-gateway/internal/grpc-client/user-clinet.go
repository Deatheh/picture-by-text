package grpcclient

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	pb "userpb"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewUserClient(serviceURL string) (*UserClient, error) {
	// Создаём соединение без блокировки
	conn, err := grpc.NewClient(serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// Запускаем соединение в фоне
	conn.Connect()

	log.Printf("Created connection to user-service at %s (connecting in background)", serviceURL)

	return &UserClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) IsHealthy() bool {
	// Проверяем состояние соединения
	state := c.conn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}

func (c *UserClient) Register(ctx context.Context, email, password string) (bool, string, error) {
	req := &pb.RegisterRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.client.Register(ctx, req)
	if err != nil {
		log.Printf("Register RPC failed: %v", err)
		return false, "", fmt.Errorf("gRPC call Register failed: %w", err)
	}

	log.Printf("Register response: success=%v, message=%s", resp.Success, resp.Uuid)
	return resp.Success, resp.Uuid, nil
}
