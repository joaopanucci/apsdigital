package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/services"
)

type TabletController struct {
	tabletService *services.TabletService
}

func NewTabletController(tabletService *services.TabletService) *TabletController {
	return &TabletController{
		tabletService: tabletService,
	}
}

type SearchAgentRequest struct {
	CPF string `form:"cpf" binding:"required"`
}

type TabletRequestRequest struct {
	Type          string `json:"type" binding:"required"`
	AgentCPF      string `json:"agent_cpf" binding:"required"`
	Reason        string `json:"reason"`
	BODocumentURL string `json:"bo_document_url"`
}

func (c *TabletController) SearchAgent(ctx *gin.Context) {
	var req SearchAgentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := c.tabletService.SearchAgentByCPF(ctx.Request.Context(), req.CPF)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agent)
}

func (c *TabletController) RequestTablet(ctx *gin.Context) {
	var req TabletRequestRequest
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

	switch req.Type {
	case "new":
		err := c.tabletService.RequestNewTablet(ctx.Request.Context(), req.AgentCPF, userEntity.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "return":
		err := c.tabletService.RequestTabletReturn(ctx.Request.Context(), req.AgentCPF)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "broken":
		err := c.tabletService.ReportTabletBroken(ctx.Request.Context(), req.AgentCPF, req.Reason)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "stolen":
		if req.BODocumentURL == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "BO document is required for stolen tablets"})
			return
		}
		err := c.tabletService.ReportTabletStolen(ctx.Request.Context(), req.AgentCPF, req.BODocumentURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request type"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tablet request submitted successfully"})
}

func (c *TabletController) GetTabletRequests(ctx *gin.Context) {
	// Get user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	// Build filters based on user role
	filters := make(map[string]interface{})

	// Non-admin users only see requests from their municipality
	if userEntity.Role.Name != "ADM" {
		filters["municipality_id"] = userEntity.MunicipalityID
	}

	requests, err := c.tabletService.GetTabletRequests(ctx.Request.Context(), filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

func (c *TabletController) ApproveRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Get user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	err = c.tabletService.ApproveTabletRequest(ctx.Request.Context(), uint(id), userEntity.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request approved successfully"})
}

func (c *TabletController) RejectRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Get user from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEntity := user.(*entities.User)

	type RejectRequest struct {
		Reason string `json:"reason" binding:"required"`
	}

	var req RejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.tabletService.RejectTabletRequest(ctx.Request.Context(), uint(id), userEntity.ID, req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request rejected successfully"})
}
