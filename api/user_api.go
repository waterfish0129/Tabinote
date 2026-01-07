package api

import (
	"github.com/gin-gonic/gin"
)

type UserApi struct {
}

func NewUserApi() *UserApi {
	return &UserApi{}
}

// @Tags 用戶管理
// @Summary 用戶登入
// @Description 用戶登入詳情描述
// @Success 200 {string} string "登入成功"
// @Failure 401 {string} string "登入失敗"
// @Router /api/v1/user/login [post]
func (m *UserApi) Login(c *gin.Context) {
	req := Build(c)
	req.OK(ResponseJson{
		Data: "",
		Code: "",
		Msg:  "",
	})
}
