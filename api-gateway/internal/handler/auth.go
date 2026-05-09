package handler

import (
	"api-gateway/internal/dpo"
	"api-gateway/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleRegistration(ctx *gin.Context) {
	if h.userClient == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "user-service is not available",
		})
		return
	}

	var input dpo.Regisration
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error getting request data: invalid request body"})
		return
	}

	if input.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error handling register: empty email"})
		return
	}
	if !pkg.CheckIsEmailAllowed(input.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error handling adding: invalid email"})
		return
	}
	if input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error handling register: empty password"})
		return
	}

	success, message, err := h.userClient.Register(
		ctx.Request.Context(),
		input.Email,
		input.Password,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Service error: " + err.Error(),
		})
		return
	}

	if !success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
