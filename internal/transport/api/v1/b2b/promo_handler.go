package b2b

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"solution/internal/shared/models/b2b/dto"
	"strconv"
)

func (h *Handler) CreatePromo(c *gin.Context) {
	var req dto.PromoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding company:", err)
		c.JSON(http.StatusBadRequest, dto.ErrorBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		log.Println("Error validating company:", err)
		c.JSON(http.StatusBadRequest, dto.ErrorBadRequest)
		return
	}

	companyID := c.GetString("company_id")

	req.CompanyID = companyID

	promoID, err := h.Promo.CreatePromo(req)
	if err != nil {
		log.Println("Error creating promo:", err)
		c.JSON(http.StatusBadRequest, dto.ErrorInternalServer)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": promoID})
}

func (h *Handler) GetPromos(c *gin.Context) {
	companyID := c.GetString("company_id")
	fmt.Println(companyID)
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	sortBy := c.Query("sort_by") // Изменено на Query
	country := c.QueryArray("country")

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

	if sortBy != "" && sortBy != "active_from" && sortBy != "active_until" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort_by parameter"})
		return
	}

	promos, totalCount, err := h.Promo.GetPromos(companyID, limit, offset, sortBy, country)
	if err != nil {
		log.Println("Error getting promos:", err)
		c.JSON(http.StatusBadRequest, dto.ErrorBadRequest.Message)
		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(totalCount, 10))

	c.JSON(http.StatusOK, promos)
}

func (h *Handler) GetPromoByID(c *gin.Context) {
	companyID := c.GetString("company_id")

	promoID := c.Param("id")

	promo, err := h.Promo.GetPromoByID(companyID, promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, dto.ErrorPromoNotFound.Error())
		} else if errors.Is(err, dto.ErrorNoAccess) {
			log.Println("Access denied to promo:", err)
			c.JSON(http.StatusForbidden, dto.ErrorNoAccessToPromo.Error())
		} else {
			log.Println("Error fetching promo:", err)
			c.JSON(http.StatusBadRequest, dto.ErrorBadRequest.Message)
		}
		return
	}

	c.JSON(http.StatusOK, promo)
}

func (h *Handler) UpdatePromo(c *gin.Context) {
	var req dto.PromoPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding promo patch request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		log.Println("Error validating promo patch request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID := c.GetString("company_id")
	promoID := c.Param("id")

	updatedPromo, err := h.Promo.UpdatePromo(companyID, promoID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, dto.ErrNotFound) {
			log.Println("Promo not found:", err)
			c.JSON(http.StatusNotFound, dto.ErrorPromoNotFound.Error())
		} else if errors.Is(err, dto.ErrorNoAccess) {
			log.Println("Access denied to promo:", err)
			c.JSON(http.StatusForbidden, dto.ErrorNoAccessToPromo.Error())
		} else if errors.Is(err, dto.ErrBadRequest) {
			log.Println("Error patching promo:", err)
			c.JSON(http.StatusBadRequest, dto.ErrorBadRequest.Message)
		} else {
			log.Println("Error fetching promo:", err)
			c.JSON(http.StatusBadRequest, dto.ErrorBadRequest.Message)
		}
		return
	}

	c.JSON(http.StatusOK, updatedPromo)
}

func (h *Handler) GetPromoStat(c *gin.Context) {
	companyID := c.GetString("company_id")
	promoID := c.Param("id")

	promoStat, err := h.Promo.GetPromoStatByID(companyID, promoID)
	if err != nil {
		if errors.Is(err, dto.ErrorPromoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": dto.ErrorPromoNotFound.Error()})
		} else if errors.Is(err, dto.ErrorNoAccessToPromo) {
			c.JSON(http.StatusForbidden, gin.H{"error": dto.ErrorNoAccessToPromo.Error()})
		} else {
			log.Println("Error getting promo stat:", err)
			c.JSON(http.StatusBadRequest, dto.ErrorBadRequest.Message)
		}
		return
	}

	c.JSON(http.StatusOK, promoStat)
}
