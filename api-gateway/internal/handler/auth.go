package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handler) HandleRegistration(ctx *gin.Context) {
	if h.userClient == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "user-service is not available",
		})
		return
	}

	var input registerInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed: " + err.Error(),
		})
		return
	}

	// Вызываем gRPC метод user-service
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
