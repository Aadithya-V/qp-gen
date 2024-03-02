package controllers

import (
	"net/http"

	"github.com/Aadithya-V/qp-gen/pkg/api/v1/models"
	"github.com/Aadithya-V/qp-gen/pkg/api/v1/services"
	"github.com/gin-gonic/gin"
)

// @Summary	Generate Q Paper Sets In Latex .tex format
// @Schemes
// @Description	Generate Q Paper Sets In Latex .tex format
// @Tags			LATEX
// @Produce		json
//
// @Param			GenerateQPaperRequestsRequest	body		models.GenerateQpaperSetsInLatexRequest	true	"request body"
// @Success		200
// @Security		BearerAuth
// @Router			/api/v1/generate-latex-q-paper-sets [post]
func GenerateQpaperSetsInLatex(c *gin.Context) {
	req := models.GenerateQpaperSetsInLatexRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "all params are required")
		return
	}

	err := services.GenerateQpaperSetsInLatex(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, "")
}
