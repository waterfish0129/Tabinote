package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/waterfish0129/Tabinote/conf"
	_ "github.com/waterfish0129/Tabinote/docs"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/middleware"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type IFnRegisterRoute = func(rg *gin.RouterGroup)

var (
	gfnRoutes []IFnRegisterRoute
)

func RegisterRoute(fn IFnRegisterRoute) {
	if fn == nil {
		return
	}
	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() {
	/*===========================================================================================*/
	//為了能優雅的退出gin應用 創建一個通道 監聽關閉事件和ctrl+c事件
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	/*===========================================================================================*/
	//定義API路由
	r := gin.Default()
	r.Use(
		//讓API允許跨域訪問
		middleware.Cors(),
		//把每個請求都加上uuid
		middleware.RequestIDMiddleware(),
		//加上錯誤捕捉  確保API在執行中不會被panic導致中斷
		middleware.RecoveryMiddleware(),
	)
	routerGroup := r.Group("/api/v1")

	/*===========================================================================================*/
	//獲取全部的API並且初始化
	InitBasePlatFormRoutes()

	/*===========================================================================================*/
	//註冊全部的API
	for _, fnRegisterRouter := range gfnRoutes {
		fnRegisterRouter(routerGroup)
	}

	/*===========================================================================================*/
	// 註冊swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	/*===========================================================================================*/
	//獲取組態檔中的端口設定
	stPort := conf.GetServerPort()
	if stPort == "" {
		stPort = "8787"
	}
	/*===========================================================================================*/
	//創建一個伺服器
	server := &http.Server{
		Addr:    ":" + stPort,
		Handler: r,
	}

	global.Logger.Info(fmt.Sprintf("Start server listen at %s", stPort))

	/*===========================================================================================*/
	//開啟一個執行續給伺服器運行  伺服器監聽錯誤   如果遇到的錯誤不是關閉   就跳錯誤訊息
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error(fmt.Sprintf("Start server listen err: %s", err))
			return
		}
	}()
	/*===========================================================================================*/
	//等待接收關閉和ctrl+c訊號  如果有收到訊號  程式才會繼續往下走
	<-ctx.Done()

	/*===========================================================================================*/
	//創建一個5秒的延遲 時間到的話就會關閉伺服器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error(fmt.Sprintf("Stop server err: %s", err))
	}

	global.Logger.Info(fmt.Sprintf("Stop server success"))
}

func InitBasePlatFormRoutes() {
	InitUserRouter()
}
