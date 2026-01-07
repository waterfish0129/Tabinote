package router

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/api"
)

func RegisterUser(rg *gin.RouterGroup) {
	userApi := api.NewUserApi()

	user := rg.Group("/user")
	{
		user.POST("/login", userApi.Login)
	}
}
