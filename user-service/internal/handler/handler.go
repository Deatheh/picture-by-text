package handler

import (
	"context"
	"log"
	"user-service/internal/config"
	"user-service/internal/service"
	"user-service/pkg"

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

func (h *UserHandler) Login(ctx context.Context, req *pb.RegisterRequest) (*pb.LoginResponse, error) {
	log.Printf("Login called: email=%s", req.Email)

	user, err := h.service.User.GetByEmail(req.Email)
	if err != nil {
		return &pb.LoginResponse{Success: false}, nil
	}
	if !h.service.User.CheckPassword(user.Password, req.Password) {
		return &pb.LoginResponse{Success: false}, nil
	}

	accessToken, err := pkg.GenerateAccessToken(user.Uuid, h.envConf.JWT.AccessTTL, h.envConf.JWT.Secret)
	if err != nil {
		return &pb.LoginResponse{Success: false}, nil
	}
	refreshToken, err := pkg.GenerateRefreshToken(user.Uuid, h.envConf.JWT.RefreshTTL/3600, h.envConf.JWT.Secret)
	if err != nil {
		return &pb.LoginResponse{Success: false}, nil
	}

	// Сохраняем refresh‑токен в Redis
	if err := h.service.User.SaveRefreshToken(ctx, user.Uuid, refreshToken, h.envConf.JWT.RefreshTTL); err != nil {
		log.Printf("Failed to save refresh token: %v", err)
	}

	return &pb.LoginResponse{
		Success:      true,
		AccsesToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// 1. Парсим refresh‑токен
	claims, err := pkg.ParseToken(req.RefreshToken, h.envConf.JWT.Secret)
	if err != nil {
		return &pb.RefreshTokenResponse{Success: false, Message: "invalid token"}, nil
	}
	if claims.Type != "refresh" {
		return &pb.RefreshTokenResponse{Success: false, Message: "invalid token type"}, nil
	}
	userID := claims.UUID

	// 2. Проверяем наличие токена в Redis
	storedToken, err := h.service.User.GetRefreshToken(ctx, userID)
	if err != nil || storedToken != req.RefreshToken {
		return &pb.RefreshTokenResponse{Success: false, Message: "token not found or revoked"}, nil
	}

	// 3. Генерируем новый access‑токен
	newAccessToken, err := pkg.GenerateAccessToken(userID, h.envConf.JWT.AccessTTL, h.envConf.JWT.Secret)
	if err != nil {
		return &pb.RefreshTokenResponse{Success: false, Message: "internal error"}, nil
	}

	return &pb.RefreshTokenResponse{Success: true, AccessToken: newAccessToken}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Определяем значения по умолчанию
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	limit := int(req.Limit)
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, total, err := h.service.User.ListUsers(page, limit)
	if err != nil {
		log.Printf("ListUsers failed: %v", err)
		return &pb.ListUsersResponse{Users: []*pb.UserItem{}, Total: 0}, nil
	}

	var pbUsers []*pb.UserItem
	for _, u := range users {
		// Определяем роль по role_id
		role := "user"
		if u.RoleID == 1 {
			role = "admin"
		}
		pbUsers = append(pbUsers, &pb.UserItem{
			Id:            u.Uuid,
			Email:         u.Email,
			Role:          role,
			RequestsCount: 0,
		})
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: int32(total),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Проверка, что не удаляем самого себя (ID из контекста)
	// Временно пропустим проверку

	err := h.service.User.DeleteUser(req.UserId)
	if err != nil {
		log.Printf("DeleteUser failed: %v", err)
		return &pb.DeleteUserResponse{Success: false, Message: err.Error()}, nil
	}

	// Удаляем refresh-токен из Redis
	_ = h.service.Cache.DeleteRefreshToken(ctx, req.UserId)

	return &pb.DeleteUserResponse{Success: true, Message: "User deleted successfully"}, nil
}

func (h *UserHandler) GetUserRole(ctx context.Context, req *pb.GetUserRoleRequest) (*pb.GetUserRoleResponse, error) {
	log.Printf("GetUserRole called: user_id=%s", req.UserId)

	user, err := h.service.User.GetByID(req.UserId)
	if err != nil {
		log.Printf("GetUserRole: user not found: %v", err)
		return &pb.GetUserRoleResponse{Role: ""}, nil
	}

	role := "user"
	if user.RoleID == 1 {
		role = "admin"
	}

	log.Printf("GetUserRole: user_id=%s, role=%s", req.UserId, role)
	return &pb.GetUserRoleResponse{Role: role}, nil
}
