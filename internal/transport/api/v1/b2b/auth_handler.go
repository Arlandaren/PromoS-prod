package b2b

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"solution/internal/service/b2b"
	"solution/internal/shared/models/b2b/dto"
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
		log.Println("Error parsing request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, companyID, err := h.Auth.RegisterCompany(dto.SignUpRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, b2b.ErrEmailAlreadyRegistered) {
			log.Println("Error parsing request:", err)
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		} else {
			log.Println("Error parsing request:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token:     token,
		CompanyID: companyID,
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

	token, err := h.Auth.AuthenticateCompany(dto.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, b2b.ErrInvalidCredentials) {
			log.Println("Error parsing request:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Println("Error parsing request:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
