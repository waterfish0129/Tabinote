package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/waterfish0129/Tabinote/global"
)

type HealthApi struct {
	*BaseApi
}

func NewHealthApi() HealthApi {
	return HealthApi{
		BaseApi: NewBaseApi(),
	}
}

// @Tags Health 伺服器狀態檢查
// @Summary 測試
// @Description
// @Success 200 {string} string "成功"
// @Failure 401 {string} string "失敗"
// @Router /api/v1/health/db [get]
func (m HealthApi) DB(c *gin.Context) {
	m.BuildRequest(c)
	var serverVersion string
	err := global.DBPool.QueryRow(m.Context(), "select version()").Scan(&serverVersion)
	if err != nil {
		m.Fail(
			ResponseJson{
				Data: err.Error(),
				Code: "",
				Msg:  "db query fail",
			})
		return
	}

	m.OK(ResponseJson{
		Data: serverVersion,
		Code: "",
		Msg:  "db query pass",
	})
}

// @Tags Health 伺服器狀態檢查
// @Summary 測試token解碼
// @Description
// @Security BearerAuth
// @Success 200 {string} string "成功"
// @Failure 401 {string} string "失敗"
// @Router /api/v1/health/parseToken [post]
func (m HealthApi) ParseToken(c *gin.Context) {
	m.BuildRequest(c)

	if !c.MustGet("isLogin").(bool) {
		m.Fail(ResponseJson{
			Msg:  "需要登入哦哦哦",
			Code: "",
			Data: nil,
		})
		return
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"userId": c.MustGet("userId"),
			"role":   c.MustGet("role"),
			"plan":   c.MustGet("plan")},
		Code: "",
		Msg:  "",
	})
}

// @Tags Health 伺服器狀態檢查
// @Summary 生成一個uuid
// @Description
// @Success 200 {string} string "成功"
// @Failure 401 {string} string "失敗"
// @Router /api/v1/health/uuid [get]
func (m HealthApi) Uuid(c *gin.Context) {
	m.BuildRequest(c)
	id, _ := uuid.NewUUID()
	m.OK(ResponseJson{
		Data: id,
		Code: "",
		Msg:  "",
	})
}
