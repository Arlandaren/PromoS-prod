package b2c

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"solution/internal/shared/models/b2c/dto"
	"strconv"
)

func (h *Handler) GetPromosForUser(c *gin.Context) {
	userID := c.GetString("user_id")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	category := c.Query("category")
	activeStr := c.Query("active")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println("Error parsing limit:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Println("Error parsing offset:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	var active *bool
	if activeStr != "" {
		activeValue, err := strconv.ParseBool(activeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'active' parameter"})
			return
		}
		active = &activeValue
	}

	promos, totalCount, err := h.Promo.GetPromosForUser(userID, limit, offset, category, active)
	if err != nil {
		log.Println("Error getting promos:", err)
		c.JSON(http.StatusBadRequest, dto.ErrBadRequest)
		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(totalCount, 10))

	c.JSON(http.StatusOK, promos)

}

func (h *Handler) GetPromo(c *gin.Context) {
	promoID := c.Param("id")
	userID := c.GetString("user_id")

	promo, err := h.Promo.GetPromo(promoID, userID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, dto.ErrNotFound)
			return
		} else {
			log.Println("Error getting promo:", err)
			c.JSON(http.StatusBadRequest, dto.ErrNotFound)
			return
		}
	}

	c.JSON(http.StatusOK, promo)
}

func (h *Handler) LikePromo(c *gin.Context) {
	promoID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.Promo.LikePromo(promoID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, dto.ErrNotFound)
			return
		} else {
			log.Println("Error liking promo:", err)
			c.JSON(http.StatusBadRequest, dto.ErrBadRequest)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) UnlikePromo(c *gin.Context) {
	promoID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.Promo.UnlikePromo(promoID, userID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, dto.ErrNotFound)
			return
		} else {
			log.Println("Error unliking promo:", err)
			c.JSON(http.StatusBadRequest, dto.ErrBadRequest)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// AddComment добавляет комментарий к промокоду
func (h *Handler) AddComment(c *gin.Context) {
	userID := c.GetString("user_id")
	promoID := c.Param("id")

	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding comment request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Проверяем длину текста комментария
	if err := req.Validate(); err != nil {
		log.Println("Error validating comment:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text length must be between 10 and 1000 characters"})

		return
	}

	comment, err := h.Promo.AddComment(userID, promoID, req.Text)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Promo is not found"})
			return
		}
		log.Println("Error adding comment:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetComments(c *gin.Context) {
	promoID := c.Param("id")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	comments, totalCount, err := h.Promo.GetComments(promoID, limit, offset)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Promo not found"})
			return
		}
		log.Println("Error getting comments:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get comments"})
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", totalCount))
	c.JSON(http.StatusOK, comments)
}

// GetComment получает конкретный комментарий по его ID
func (h *Handler) GetComment(c *gin.Context) {
	promoID := c.Param("id")
	commentID := c.Param("comment_id")

	comment, err := h.Promo.GetComment(promoID, commentID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Comment not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		log.Println("Error getting comment:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// EditComment редактирует текст комментария
func (h *Handler) EditComment(c *gin.Context) {
	userID := c.GetString("user_id")
	promoID := c.Param("id")
	commentID := c.Param("comment_id")

	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding comment request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	comment, err := h.Promo.EditComment(userID, promoID, commentID, req.Text)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Comment not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if errors.Is(err, dto.ErrNoAccess) {
			log.Println("No access to edit this comment:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "No access to edit this comment"})
			return
		}
		log.Println("Error editing comment:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to edit comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment удаляет комментарий
func (h *Handler) DeleteComment(c *gin.Context) {
	userID := c.GetString("user_id")
	promoID := c.Param("id")
	commentID := c.Param("comment_id")

	err := h.Promo.DeleteComment(userID, promoID, commentID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			log.Println("Comment not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if errors.Is(err, dto.ErrNoAccess) {
			log.Println("No access to delete this comment:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "No access to delete this comment"})
			return
		}
		log.Println("Error deleting comment:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
