package handler

import (
	"api-gateway/internal/config"
	grpcclient "api-gateway/internal/grpc-client"
	"api-gateway/internal/logger"

	"github.com/gin-gonic/gin"
)

const (
	Route = "/"

	AuthRegisterRoute = "/auth/register"
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
	}

	return r
}
