package handler

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"

	"github.com/gin-gonic/gin"
)

const (
	Route = "/"
)

type Handler struct {
	envConf *config.Config
}

func NewRouter(envConf *config.Config) *Handler {
	return &Handler{envConf: envConf}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	auth := r.Group("")
	{

		auth.Use(logger.RequestLogger("auth"))
	}

	return r
}
