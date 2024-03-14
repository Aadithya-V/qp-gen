package controllers

import (
	"fmt"
	"log"
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

	// err := services.GenerateQpaperSetsInLatex(c, &req)
	/* if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} */

	c.JSON(200, "")
}

// @Summary	Generate Q Paper Sets In Latex .tex format from db
// @Schemes
// @Description	Generate Q Paper Sets In Latex .tex format from db
// @Tags			DB
// @Produce		json
//
// @Param			GenerateQpaperSetsFromDBRequest	body		models.GenerateQpaperSetsFromDBRequest{}	true	"request body"
// @Success		200
// @Security		BearerAuth
// @Router			/api/v1/q-paper-from-db [post]
func GenerateQpaperSetsFromDB(c *gin.Context) {
	req := models.GenerateQpaperSetsFromDBRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "all params are required")
		return
	}

	err := services.GenerateQpaperSetsFromDB(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, "")
}

// @Summary	Upload Question Bank in CSV format
// @Schemes
// @Description	Upload Question Bank in CSV format
// @Tags		Upload
// @Accept mpfd
// @Param academic_year path string true "academic year"
// @Param subject_code path string true "subject code"
// @Param file formData file true "File to upload"
// @Produce		json
//
// @Success		200
// @Security		BearerAuth
// @Router			/api/v1/upload/{academic_year}/{subject_code} [post]
func UploadCSV(c *gin.Context) {

	subCode := c.Param("subject_code")
	year := c.Param("academic_year")

	if len(subCode) == 0 || len(year) == 0 {
		err := fmt.Errorf("required path params subject_code and academic_year")
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, err.Error())
		c.Abort()
		return
	}

	if file == nil {
		log.Println("no file uploaded")
		c.JSON(http.StatusBadRequest, "required file upload")
		c.Abort()
		return
	}

	log.Println(file.Filename)

	err = services.ParseAndSaveQuestionsFromCSV(c, file, subCode, year)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, err.Error())
		c.Abort()
		return
	}

	c.JSON(200, "file uploaded")
}
