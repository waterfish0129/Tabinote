package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/waterfish0129/Tabinote/api"
	consts "github.com/waterfish0129/Tabinote/global/constants"
	"github.com/waterfish0129/Tabinote/utils"
	"strings"
)

func reject() api.ResponseJson {
	return api.ResponseJson{
		Status: 401,
		Code:   consts.ErrTokenInvalid,
		Msg:    consts.MsgTokenInvalid,
		Data:   nil,
	}
}

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("isLogin", false)
		token := c.GetHeader("Authorization")
		//如果沒有帶token 他沒登入直接通過
		if token == "" {
			c.Next()
			return
		}

		//如果token格式不對 拒絕
		if !strings.HasPrefix(token, "Bearer") {
			api.Fail(c, reject())
			c.Abort()
			return
		}

		/*===========================================================================================*/
		//token解析
		//取出Token主體
		token = token[len("Bearer "):]
		//解析token  拿到自定義聲明
		myJwtClaims, err := utils.ParseToken(token)

		IsTokenExpired := errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet)
		if err != nil && IsTokenExpired {
			//如果是token過期  就回傳過期
			api.Fail(c, api.ResponseJson{
				Status: 401,
				Code:   consts.ErrTokenExpired,
				Msg:    consts.MsgTokenExpired,
				Data:   nil,
			})
			c.Abort()
			return
		}

		if err != nil {
			//如果token無效或解析有錯就拒絕
			api.Fail(c, reject())
			c.Abort()
			return
		}

		c.Set("userId", myJwtClaims.UserId)
		c.Set("role", myJwtClaims.Role)
		c.Set("plan", myJwtClaims.Plan)
		c.Set("isLogin", true)

		c.Next()
	}

}
func AuthOnlyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		isLogin, ok := c.Get("isLogin")
		if !ok || isLogin != true {
			api.Fail(c, reject())
			c.Abort()
			return
		}

		c.Next()
	}
}
