package router

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/api"
	"github.com/waterfish0129/Tabinote/middleware"
)

func RegisterAuth(rg *gin.RouterGroup) {
	authApi := api.NewAuthApi()

	authRoute := rg.Group("/auth")
	{
		authRoute.POST("/google", authApi.Google)
	}
	authRouteAuthOnly := authRoute.Use(middleware.AuthOnlyMiddleware())
	{
		authRouteAuthOnly.GET("/me", authApi.Me)
	}
}
