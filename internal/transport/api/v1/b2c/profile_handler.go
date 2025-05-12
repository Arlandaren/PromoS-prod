package b2c

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"solution/internal/shared/models/b2c/dto"
)

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	fmt.Println(userID)

	profile, err := h.Profile.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Profile.UpdateProfile(userID, req)
	if err != nil {
		if !errors.Is(err, dto.ErrNoFieldsToUpdate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	currentProfile, err := h.Profile.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, currentProfile)
}
