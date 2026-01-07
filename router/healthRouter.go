package router

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/api"
	"github.com/waterfish0129/Tabinote/middleware"
)

func RegisterHealth(rg *gin.RouterGroup) {
	healthApi := api.NewHealthApi()
	rgUser := rg.Group("health")
	{
		rgUser.GET("/db", healthApi.DB)
		rgUser.GET("/uuid", healthApi.Uuid)
		rgUser.Use(middleware.AuthOnlyMiddleware()).POST("/parseToken", healthApi.ParseToken)
	}
}
