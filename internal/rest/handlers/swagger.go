package handlers

import (
	_ "github.com/HunterGooD/go_test_task/docs"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewSwaggerHandler(r *gin.Engine, logger interfaces.Logger) {
	logger.Info("swagger create handler")
	// swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	logger.Info("SWAGGER register handler path GET `/swagger/*any`")
}
