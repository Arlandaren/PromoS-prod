package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	repo "solution/internal/repository/b2c"
	"solution/internal/service/services"
	"solution/internal/shared/models/b2b/dto"
	"solution/internal/shared/utils"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorUnauthorized)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, dto.ErrorUnauthorized)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorUnauthorized)
			c.Abort()
			return
		}

		userId := claims.UserID
		if userId == "" {
			fmt.Println("Invalid user_id in claims")
			c.JSON(http.StatusUnauthorized, dto.ErrorUnauthorized)
			c.Abort()
			return
		}

		var repository repo.AuthRepository

		err = services.GetService(&repository)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		existingToken, err := repository.ValidateToken(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorUnauthorized)
			return
		}

		if existingToken != tokenString {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorUnauthorized)
			return
		}

		c.Set("user_id", userId)

		c.Next()
	}
}
