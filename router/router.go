package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/waterfish0129/Tabinote/docs"
	"github.com/waterfish0129/Tabinote/middleware"
	"time"
)

type Router struct {
	engine *gin.Engine
}

func InitRouter() *Router {
	r := gin.Default()
	/*===========================================================================================*/
	// 註冊swagger 放在 routerGroup 之外。因為 Swagger 通常不需要經過 Cors、RequestID 等所有 API 中間層的處理。
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	/*===========================================================================================*/
	//加入中間層
	r.Use(
		//加上錯誤捕捉  確保API在執行中不會被panic導致中斷
		middleware.RecoveryMiddleware(),
		//讓API允許跨域訪問
		middleware.Cors(),
		//把每個請求都加上uuid
		middleware.RequestIDMiddleware(),
		//加上timeout
		middleware.TimeoutMiddleware(viper.GetDuration("api.timeout_s")*time.Second),
		//加上身分認證檢查
		middleware.AuthMiddleware(),
	)

	return &Router{
		engine: r,
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) RegisterAPIs() {
	api := r.engine.Group("/api/v1")

	RegisterUser(api)
	RegisterHealth(api)
	RegisterAuth(api)
}
