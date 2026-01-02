package cmd

import (
	"github.com/waterfish0129/Tabinote/conf"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/router"
)

func Start() {
	/*===========================================================================================*/
	//初始化系統配置文件
	conf.IntiConfig()

	/*===========================================================================================*/
	//初始化日誌
	global.Logger = conf.InitLogger()

	/*===========================================================================================*/
	//初始化路由
	router.InitRouter()
}

func Clean() {}
