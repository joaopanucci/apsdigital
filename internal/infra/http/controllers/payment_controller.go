package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/services"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

type CreatePaymentRequest struct {
	FileURL        string `json:"file_url" binding:"required"`
	Competence     string `json:"competence" binding:"required"`
	MunicipalityID uint   `json:"municipality_id" binding:"required"`
}

type PaymentFilters struct {
	Year           string `form:"year"`
	Month          string `form:"month"`
	Competence     string `form:"competence"`
	MunicipalityID string `form:"municipality_id"`
	Page           int    `form:"page"`
	Limit          int    `form:"limit"`
}

func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (set by auth middleware)
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	payment := &entities.Payment{
		FileURL:        req.FileURL,
		Competence:     req.Competence,
		UploadedBy:     userEntity.ID,
		MunicipalityID: req.MunicipalityID,
	}

	if err := c.paymentService.CreatePayment(ctx.Request.Context(), payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, payment)
}

func (c *PaymentController) GetPayments(ctx *gin.Context) {
	var filters PaymentFilters
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	// Build filters map
	filterMap := make(map[string]interface{})

	// Add municipality filter for non-admin users
	if userEntity.Role.Name != "ADM" {
		filterMap["municipality_id"] = userEntity.MunicipalityID
	} else if filters.MunicipalityID != "" {
		municipalityID, err := strconv.ParseUint(filters.MunicipalityID, 10, 32)
		if err == nil {
			filterMap["municipality_id"] = uint(municipalityID)
		}
	}

	if filters.Year != "" {
		year, err := strconv.Atoi(filters.Year)
		if err == nil {
			filterMap["year"] = year
		}
	}

	if filters.Month != "" {
		month, err := strconv.Atoi(filters.Month)
		if err == nil {
			filterMap["month"] = month
		}
	}

	if filters.Competence != "" {
		filterMap["competence"] = filters.Competence
	}

	var payments []*entities.Payment
	var err error

	if filters.Page > 0 && filters.Limit > 0 {
		payments, err = c.paymentService.GetPaymentsWithPagination(ctx.Request.Context(), filterMap, filters.Page, filters.Limit)
	} else {
		payments, err = c.paymentService.GetAllPayments(ctx.Request.Context(), filterMap)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payments)
}

func (c *PaymentController) GetPaymentByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := c.paymentService.GetPaymentByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Check access permissions
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	// Admin can access all, others only their municipality
	if userEntity.Role.Name != "ADM" && payment.MunicipalityID != userEntity.MunicipalityID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

func (c *PaymentController) UpdatePayment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := &entities.Payment{
		ID:         uint(id),
		FileURL:    req.FileURL,
		Competence: req.Competence,
	}

	if err := c.paymentService.UpdatePayment(ctx.Request.Context(), payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

func (c *PaymentController) DeletePayment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := c.paymentService.DeletePayment(ctx.Request.Context(), uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}

func (c *PaymentController) GetCompetences(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	var municipalityID *uint
	if userEntity.Role.Name != "ADM" {
		municipalityID = &userEntity.MunicipalityID
	}

	competences, err := c.paymentService.GetCompetences(ctx.Request.Context(), municipalityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"competences": competences})
}

func (c *PaymentController) GetYears(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	var municipalityID *uint
	if userEntity.Role.Name != "ADM" {
		municipalityID = &userEntity.MunicipalityID
	}

	years, err := c.paymentService.GetYears(ctx.Request.Context(), municipalityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"years": years})
}

func (c *PaymentController) ViewPDF(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := c.paymentService.GetPaymentByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Check access permissions
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	// Admin can access all, others only their municipality
	if userEntity.Role.Name != "ADM" && payment.MunicipalityID != userEntity.MunicipalityID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Return the PDF URL for frontend to handle
	ctx.JSON(http.StatusOK, gin.H{
		"file_url": payment.FileURL,
		"filename": "pagamento-" + payment.Competence + ".pdf",
	})
}

func (c *PaymentController) DownloadPDF(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := c.paymentService.GetPaymentByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Check access permissions
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	// Admin can access all, others only their municipality
	if userEntity.Role.Name != "ADM" && payment.MunicipalityID != userEntity.MunicipalityID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Return the PDF URL for frontend to handle download
	ctx.JSON(http.StatusOK, gin.H{
		"file_url": payment.FileURL,
		"filename": "pagamento-" + payment.Competence + ".pdf",
	})
}