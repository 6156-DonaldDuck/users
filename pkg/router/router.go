package router

import (
	"github.com/6156-DonaldDuck/users/docs"
	"github.com/6156-DonaldDuck/users/pkg/config"
	"github.com/6156-DonaldDuck/users/pkg/router/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) // use customized cors middleware
	r.Use(middleware.Security())
	docs.SwaggerInfo.BasePath = config.Configuration.Mysql.Host

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	InitUserRouters(r)
	InitAuthRouters(r)

	r.Run(":" + config.Configuration.Port)
}