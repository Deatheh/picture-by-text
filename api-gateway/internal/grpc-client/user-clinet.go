package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	if err := waitForConnection(ctx, conn); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	log.Printf("Connected to user-service at %s", serviceURL)

	return &UserClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func waitForConnection(ctx context.Context, conn *grpc.ClientConn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			state := conn.GetState()
			if state == connectivity.Ready {
				return nil
			}
			if !conn.WaitForStateChange(ctx, state) {
				return fmt.Errorf("connection state change failed")
			}
		}
	}
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) IsHealthy() bool {
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
		return false, "", fmt.Errorf("gRPC call Register failed: %w", err)
	}

	return resp.Success, resp.Message, nil
}
