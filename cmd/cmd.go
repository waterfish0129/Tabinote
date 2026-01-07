package cmd

import (
	"context"
	"fmt"
	"github.com/waterfish0129/Tabinote/conf"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/router"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	/*===========================================================================================*/
	//初始化系統配置文件
	conf.IntiConfig()

	/*===========================================================================================*/
	//初始化日誌
	global.Logger = conf.InitLogger()

	/*===========================================================================================*/
	//初始化資料庫連線
	if err := conf.InitPostgres(); err != nil {
		global.Logger.Fatal("database connection failed", zap.Error(err))
	}
	global.DBPool = conf.Pool
	/*===========================================================================================*/
	//初始化路由
	r := router.InitRouter()
	//註冊API
	r.RegisterAPIs()

	/*===========================================================================================*/
	//設定伺服器
	stPort := conf.GetServerPort()
	server := &http.Server{
		Addr:    ":" + stPort,
		Handler: r.Engine(),
	}

	/*===========================================================================================*/
	//為了能優雅的退出gin應用 創建一個通道 監聽關閉事件和ctrl+c事件
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	/*===========================================================================================*/
	//開啟一個執行續給伺服器運行  伺服器監聽錯誤   如果遇到的錯誤不是關閉   就跳錯誤訊息
	go func() {
		global.Logger.Info(fmt.Sprintf("Start server listen at %s", stPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error("Start server listen err", zap.Error(err))
			return
		}
	}()
	/*===========================================================================================*/
	//等待接收關閉和ctrl+c訊號  如果有收到訊號  程式才會繼續往下走
	<-ctx.Done()

	/*===========================================================================================*/
	//創建一個5秒的延遲 時間到的話就會關閉伺服器
	global.Logger.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		global.Logger.Error("Stop server err", zap.Error(err))
	}

	global.Logger.Info("Stop server success")
}

func Clean() {
	//關閉數據庫
	if global.DBPool != nil {
		global.DBPool.Close()
	}
}
