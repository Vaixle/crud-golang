package v1

import (
	"github.com/Vaixle/crud-golang/internal/controller/http/midleware"
	"github.com/Vaixle/crud-golang/internal/entity"
	"github.com/Vaixle/crud-golang/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// Swagger docs.
	_ "github.com/Vaixle/crud-golang/docs"
)

// NewRouter -
// Swagger spec:

// @title           GOLANG CRUD
// @version         1.0
// @description     API for TODO tasks
// @termsOfService  http://swagger.io/terms/

// @contact.name   Petr Petushkov
// @contact.url    https://t.me/vaixle

// @license.name  MIT
// @license.url   https://github.com/Vaixle/empha-soft/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func NewRouter(handler *gin.Engine, useCase entity.TodoUseCase, l logger.Interface) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers
	h := handler.Group("/api/v1")
	h.Use(midleware.BasicAuth())
	{
		newTODORoutes(h, useCase, l)
	}
}
