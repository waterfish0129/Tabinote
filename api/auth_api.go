package api

import (
	"github.com/gin-gonic/gin"
	"github.com/waterfish0129/Tabinote/dao"
	"github.com/waterfish0129/Tabinote/global"
	consts "github.com/waterfish0129/Tabinote/global/constants"
	"github.com/waterfish0129/Tabinote/service"
	"github.com/waterfish0129/Tabinote/service/dto"
	"go.uber.org/zap"
)

type AuthApi struct {
	authService *service.AuthService
}

func NewAuthApi() *AuthApi {
	return &AuthApi{
		authService: service.NewAuthService(dao.NewUserDao(), dao.NewPlanDao()),
	}
}

// @Tags Auth
// @Summary google登入
// @Description
// @Param data body dto.GoogleLoginDTO true "payload"
// @Success 200 {string} string "成功"
// @Failure 401 {string} string "失敗"
// @Router /api/v1/auth/google [post]
func (m *AuthApi) Google(c *gin.Context) {
	/*===========================================================================================*/
	//創建API請求並綁定數據
	req := Build(c)
	var iGoogleLoginDTO dto.GoogleLoginDTO
	if req.Bind(BindDataOption{DTO: &iGoogleLoginDTO}).HasError() {
		return
	}

	googleLoginResponse, err := m.authService.GoogleLogin(req.Context(), iGoogleLoginDTO.IdToken)

	if err != nil {
		global.Logger.Warn("google login error",
			zap.Error(err),
			zap.String("request_id", req.RequestID()),
			zap.String("location", "auth_api.Google"))
		req.Fail(ResponseJson{
			Code: consts.ErrLoginFailed,
			Msg:  consts.MsgLoginFailed,
			Data: nil,
		})
		return
	}

	req.OK(ResponseJson{
		Data: googleLoginResponse,
		Msg:  "",
		Code: "",
	})
}

// @Tags Auth
// @Summary Me
// @Description get my userInfo
// @Security BearerAuth
// @Success 200 {string} string "成功"
// @Failure 401 {string} string "失敗"
// @Router /api/v1/auth/me [get]
func (m *AuthApi) Me(c *gin.Context) {
	/*===========================================================================================*/
	//創建API請求並綁定數據
	req := Build(c)

	_, userId := req.GetUserId()
	userProfile, err := m.authService.Me(req.Context(), userId)

	if err != nil {
		global.Logger.Warn("get me error",
			zap.Error(err),
			zap.String("request_id", req.RequestID()),
			zap.String("location", "auth_api.Me"))
		req.Fail(ResponseJson{
			Code: consts.ErrGetMeFailed,
			Msg:  consts.MsgGetMeFailed,
			Data: nil,
		})
		return
	}

	req.OK(ResponseJson{
		Data: userProfile,
		Msg:  "",
		Code: "",
	})
}
