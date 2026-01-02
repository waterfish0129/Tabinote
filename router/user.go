package router

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/api"
	"net/http"
)

func InitUserRouter() {
	RegisterRoute(func(rg *gin.RouterGroup) {
		userApi := api.NewUserApi()
		rgUser := rg.Group("user")
		{
			rgUser.POST("/login", userApi.Login)
			rgUser.GET("", func(ctx *gin.Context) {
				ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"msg":  "ok",
					"data": []map[string]any{{"id": 1, "name": "水魚"}, {"id": 2, "name": "鯉魚王"}},
				})
			})
			rgUser.GET("/:id", func(ctx *gin.Context) {
				ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"msg":  "ok",
					"data": map[string]any{"id": 1, "name": "水魚"},
				})
			})
		}

	})
}
