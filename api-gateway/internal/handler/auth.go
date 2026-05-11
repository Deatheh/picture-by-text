package handler

import (
	"api-gateway/internal/dpo"
	"api-gateway/internal/entities"
	"api-gateway/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleRegistration(ctx *gin.Context) {
	if h.userClient == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "user-service is not available"})
		return
	}

	var input dpo.Regisration
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if input.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty email"})
		return
	}
	if !pkg.CheckIsEmailAllowed(input.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	if input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty password"})
		return
	}

	success, userID, err := h.userClient.Register(ctx.Request.Context(), input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Service error: " + err.Error()})
		return
	}

	if !success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Registration failed (email may already exist)"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (h *Handler) HandleLogin(c *gin.Context) {
	if h.userClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "user-service unavailable"})
		return
	}

	var input dpo.Regisration
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	success, accessToken, refreshToken, err := h.userClient.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login service error"})
		return
	}
	if !success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, &entities.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) HandleRefreshToken(c *gin.Context) {
	var req dpo.Refresh
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	success, newAccessToken, err := h.userClient.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}
	if !success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
