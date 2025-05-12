package b2b

import (
	"github.com/gin-gonic/gin"
	"log"
	"solution/internal/service/b2b"
	"solution/internal/service/services"
	"solution/internal/transport/api/v1/b2b/middleware"
)

type BusinessHandler interface {
	Route(r *gin.Engine)
	RouteBusinessAuth(r *gin.Engine)
	RouteBusinessPromo(r *gin.Engine)
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	CreatePromo(c *gin.Context)
	GetPromos(c *gin.Context)
	GetPromoByID(c *gin.Context)
	UpdatePromo(c *gin.Context)
	GetPromoStat(c *gin.Context)
}

type Handler struct {
	Auth  b2b.AuthService
	Promo b2b.PromoService
}

func NewHandler() *Handler {
	h := &Handler{}

	err := services.GetService(&h.Auth)
	if err != nil {
		log.Fatalf("Failed to get AuthService: %v", err)
	}

	err = services.GetService(&h.Promo)
	if err != nil {
		log.Fatalf("Failed to get AuthService: %v", err)
	}

	return h
}

func (h *Handler) Route(r *gin.Engine) {
	h.RouteBusinessAuth(r)
	h.RouteBusinessPromo(r)
}

func (h *Handler) RouteBusinessAuth(r *gin.Engine) {
	businessAuth := r.Group("api/business/auth")
	{
		businessAuth.POST("/sign-up", h.SignUp)
		businessAuth.POST("/sign-in", h.SignIn)
	}
}

func (h *Handler) RouteBusinessPromo(r *gin.Engine) {

	businessPromo := r.Group("api/business/promo")
	businessPromo.Use(middleware.AuthMiddleware())
	{
		businessPromo.POST("", h.CreatePromo)
		businessPromo.GET("", h.GetPromos)
		businessPromo.GET("/:id", h.GetPromoByID)
		businessPromo.PATCH("/:id", h.UpdatePromo)
		businessPromo.GET("/:id/stat", h.GetPromoStat)
	}
}
