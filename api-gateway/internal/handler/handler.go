package handler

import (
	"api-gateway/internal/config"
	grpcclient "api-gateway/internal/grpc-client"
	"api-gateway/internal/logger"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

const (
	Route = "/"

	AuthRegisterRoute = "/auth/register"
	AuthLoginRoute    = "/auth/login"
	AuthRefreshhRoute = "/auth/refresh_token"

	AdminUsersRoute = "/admin/users"
)

type Handler struct {
	envConf    *config.Config
	userClient *grpcclient.UserClient
}

func NewRouter(envConf *config.Config, userClient *grpcclient.UserClient) *Handler {
	return &Handler{envConf: envConf, userClient: userClient}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	auth := r.Group("")
	{
		auth.Use(logger.RequestLogger("auth"))
		auth.POST(AuthRegisterRoute, h.HandleRegistration)
		auth.POST(AuthLoginRoute, h.HandleLogin)
		auth.POST(AuthRefreshhRoute, h.HandleRefreshToken)
	}
	admin := r.Group("")
	admin.Use(middleware.JWTAuthMiddleware(h.envConf.JWT.Secret))
	admin.Use(middleware.AdminMiddleware(h.userClient))
	{
		admin.GET(AdminUsersRoute, h.HandleListUsers)
		admin.DELETE(AdminUsersRoute, h.HandleDeleteUser)
	}

	return r
}
