package dto

type PromoForUser struct {
	PromoID           string `json:"promo_id" binding:"required"`
	CompanyID         string `json:"company_id" binding:"required"`
	CompanyName       string `json:"company_name" binding:"required"`
	Description       string `json:"description" binding:"required"`
	ImageURL          string `json:"image_url"`
	Active            bool   `json:"active" binding:"required"`
	IsActivatedByUser bool   `json:"is_activated_by_user" binding:"required"`
	LikeCount         int    `json:"like_count" binding:"required"`
	IsLikedByUser     bool   `json:"is_liked_by_user" binding:"required"`
	CommentCount      int64  `json:"comment_count" binding:"required"`
}
