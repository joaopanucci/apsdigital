package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaopanucci/apsdigital/internal/domain/services"
)

type MunicipalityController struct {
	municipalityService *services.MunicipalityService
}

func NewMunicipalityController(municipalityService *services.MunicipalityService) *MunicipalityController {
	return &MunicipalityController{
		municipalityService: municipalityService,
	}
}

func (c *MunicipalityController) GetMunicipalities(ctx *gin.Context) {
	municipalities, err := c.municipalityService.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, municipalities)
}
