package routers

import (
	"github.com/Aadithya-V/qp-gen/docs"
	"github.com/Aadithya-V/qp-gen/pkg/api/v1/controllers"

	"github.com/gin-contrib/cors"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()

	router.Use(cors.Default())
	router.Use(gin.Recovery())

	// router.Use(middleware.RequestId())

	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/qp-gen"
	docs.SwaggerInfo.InfoInstanceName = "qp gen service"
	docs.SwaggerInfo.Title = "QP Generator"
	docs.SwaggerInfo.Description = "This Service Is The Primary Backend Endpoint"
	//docs.SwaggerInfo.Schemes = []string{"http", "https"}
	baseGroup := router.Group("qp-gen")
	{
		apiGroup := baseGroup.Group("api")
		{
			versionOne := apiGroup.Group("v1")
			// versionOne.Use(middleware.IsAuthorised())
			{

				versionOne.POST("generate-latex-q-paper-sets", controllers.GenerateQpaperSetsInLatex)

				versionOne.POST("q-paper-from-db", controllers.GenerateQpaperSetsFromDB)

				versionOne.POST("upload/:academic_year/:subject_code", controllers.UploadCSV)

			}
			apiGroup.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(1), ginSwagger.DocExpansion("none"), ginSwagger.PersistAuthorization(true)))

		}
	}
	return router
}
