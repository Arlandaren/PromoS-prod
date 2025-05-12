package b2c

import (
	"github.com/gin-gonic/gin"
	"log"
	"solution/internal/service/b2c"
	"solution/internal/service/services"
	"solution/internal/transport/api/v1/b2c/middleware"
)

type UserHandler interface {
	Route(r *gin.Engine)
	RouteUserAuth(r *gin.Engine)
}

type Handler struct {
	Auth    b2c.AuthService
	Profile b2c.ProfileService
	Promo   b2c.PromoService
}

func NewHandler() *Handler {
	h := &Handler{}

	err := services.GetService(&h.Auth)
	if err != nil {
		log.Fatalf("Failed to get AuthService: %v", err)
	}

	err = services.GetService(&h.Profile)
	if err != nil {
		log.Fatalf("Failed to get ProfileService: %v", err)
	}

	err = services.GetService(&h.Promo)
	if err != nil {
		log.Fatalf("Failed to get PromoService: %v", err)
	}

	return h
}

func (h *Handler) Route(r *gin.Engine) {
	h.RouteUserAuth(r)
	h.RouteUserProfile(r)
	h.RouteUserPromo(r)
}

func (h *Handler) RouteUserAuth(r *gin.Engine) {
	userAuth := r.Group("api/user/auth")
	{
		userAuth.POST("/sign-up", h.SignUp)
		userAuth.POST("/sign-in", h.SignIn)
	}
}

func (h *Handler) RouteUserProfile(r *gin.Engine) {
	{
		userProfile := r.Group("api/user/profile")
		userProfile.Use(middleware.AuthMiddleware())
		{
			userProfile.GET("", h.GetProfile)
			userProfile.PATCH("", h.UpdateProfile)
		}
	}
}

func (h *Handler) RouteUserPromo(r *gin.Engine) {
	{
		user := r.Group("api/user")
		user.Use(middleware.AuthMiddleware())

		user.GET("/feed", h.GetPromosForUser)

		promo := user.Group("/promo")
		promo.Use(middleware.AuthMiddleware())
		{
			promo.GET(":id", h.GetPromo)
			promo.POST(":id/like", h.LikePromo)
			promo.DELETE(":id/like", h.UnlikePromo)
			comments := promo.Group(":id/comments")
			comments.Use(middleware.AuthMiddleware())
			{
				comments.POST("", h.AddComment)
				comments.GET("", h.GetComments)

				comment := comments.Group(":comment_id")
				comment.Use(middleware.AuthMiddleware())
				{
					comment.GET("", h.GetComment)
					comment.PUT("", h.EditComment)
					comment.DELETE("", h.DeleteComment)
				}
			}
		}

	}
}
