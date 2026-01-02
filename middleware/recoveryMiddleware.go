package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/api"
	"github.com/waterfish0129/Tabinote/global"
	"go.uber.org/zap"
)

const (
	ErrInternal = "INTERNAL_ERROR"
)
const (
	MsgInternalError = "server.internal_error"
)

func RecoveryMiddleware() gin.HandlerFunc {
	/*===========================================================================================*/
	//用來保護所有的請求  如果發生不預期的錯誤  不會導致整個伺服器崩壞
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				/*===========================================================================================*/
				//把錯誤寫到日誌
				global.Logger.Error(
					"panic recovered",
					zap.String("request_id", c.GetString("request_id")),
					zap.Any("panic", r), //發生錯誤的代碼原因
					zap.Stack("stack"),  //發生錯誤的代碼位置
				)

				api.ServerFail(c, api.ResponseJson{
					Code: ErrInternal,
					Msg:  MsgInternalError,
				})
			}
		}()

		c.Next() // ⚠️ 一定要呼叫
	}
}
