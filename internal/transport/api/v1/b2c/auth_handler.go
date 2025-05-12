package b2c

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"solution/internal/service/b2c"
	"solution/internal/shared/models/b2c/dto"
	"strings"
)

func (h *Handler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(req.Email)

	if err := req.Validate(); err != nil {
		log.Println("Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, userID, err := h.Auth.RegisterUser(req)
	if err != nil {
		if errors.Is(err, b2c.ErrEmailAlreadyRegistered) {
			log.Println("Email already registered:", req.Email)
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		} else {
			log.Println("Internal server error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad req"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user_id": userID,
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(req.Email)

	if err := req.Validate(); err != nil {
		log.Println("Error parsing request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Auth.AuthenticateUser(req)
	if err != nil {
		if errors.Is(err, dto.ErrInvalidCredentials) {
			log.Println("Error parsing request:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Println("Error parsing request:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest."})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
