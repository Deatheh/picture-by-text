package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleRegistration(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Registration endpoint - user-service not yet implemented",
		"status":  "placeholder",
	})
}
