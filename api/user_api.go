package api

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/service/dto"
	"github.com/waterfish0129/Tabinote/utils"
)

type UserApi struct {
	BaseApi
}

func NewUserApi() UserApi {
	return UserApi{
		BaseApi: NewBaseApi(),
	}
}

// @Tags 用戶管理
// @Summary 用戶登入
// @Description 用戶登入詳情描述
// @Param data body dto.UserLoginDTO true "payload"
// @Success 200 {string} string "登入成功"
// @Failure 401 {string} string "登入失敗"
// @Router /api/v1/user/login [post]
func (m UserApi) Login(c *gin.Context) {
	var iUserLoginDTO dto.UserLoginDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserLoginDTO}).GetError(); err != nil {
		return
	}
	token, _ := utils.GenerateToken(10, "admin", "free")
	m.OK(ResponseJson{
		Data: token,
		Code: "",
		Msg:  "",
	})
}
