package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/services"
)

type ProfessionController struct {
	professionService *services.ProfessionService
}

func NewProfessionController(professionService *services.ProfessionService) *ProfessionController {
	return &ProfessionController{
		professionService: professionService,
	}
}

type CreateProfessionRequest struct {
	Name string `json:"name" binding:"required"`
}

func (c *ProfessionController) CreateProfession(ctx *gin.Context) {
	var req CreateProfessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profession := &entities.Profession{
		Name: req.Name,
	}

	if err := c.professionService.CreateProfession(ctx.Request.Context(), profession); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, profession)
}

func (c *ProfessionController) GetProfessions(ctx *gin.Context) {
	professions, err := c.professionService.GetAllProfessions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, professions)
}

func (c *ProfessionController) GetProfessionByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profession ID"})
		return
	}

	profession, err := c.professionService.GetProfessionByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Profession not found"})
		return
	}

	ctx.JSON(http.StatusOK, profession)
}

func (c *ProfessionController) UpdateProfession(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profession ID"})
		return
	}

	var req CreateProfessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profession := &entities.Profession{
		ID:   id,
		Name: req.Name,
	}

	if err := c.professionService.UpdateProfession(ctx.Request.Context(), profession); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profession)
}

func (c *ProfessionController) DeleteProfession(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profession ID"})
		return
	}

	if err := c.professionService.DeleteProfession(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profession deleted successfully"})
}
