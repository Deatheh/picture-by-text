package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleListUsers(c *gin.Context) {
	if h.userClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "user-service unavailable"})
		return
	}

	// Парсим параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.userClient.ListUsers(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
}

func (h *Handler) HandleDeleteUser(c *gin.Context) {
	if h.userClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "user-service unavailable"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "error parsing user data: id not found"})
	}

	success, message, err := h.userClient.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
