package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	handlers2 "payment/internal/http/handlers"
	"payment/internal/repositories"
	pgGorm "payment/internal/repositories/pg-gorm"
	"payment/internal/services"
)

func ApplicationV1Router(
	newPgRepo pgGorm.PGInterface,
	router *gin.Engine,
) {
	routerV1 := router.Group("/v1")
	{
		// Swagger
		routerV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Migrations
		MigrateRoutes(routerV1, handlers2.NewMigrationHandler(newPgRepo))

		// Auth for User
		authUserRepo := repositories.NewAuthUserRepository(newPgRepo)
		authUserService := services.NewAuthUserService(authUserRepo, newPgRepo)
		AuthorizationUserRoutes(routerV1, handlers2.NewAuthUserHandler(newPgRepo, authUserService))
	}
}

func MigrateRoutes(router *gin.RouterGroup, handler *handlers2.MigrationHandler) {
	routerAuth := router.Group("/internal")
	{
		routerAuth.POST("/migrate", handler.Migrate)
	}
}

func AuthorizationUserRoutes(router *gin.RouterGroup, handler *handlers2.AuthUserHandler) {
	routerAuth := router.Group("/auth")
	{
		routerAuth.POST("/login", handler.Login)
		routerAuth.POST("/register", handler.Register)
	}
}
