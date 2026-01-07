package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/waterfish0129/Tabinote/global"
	consts "github.com/waterfish0129/Tabinote/global/constants"
	"go.uber.org/zap"
	"net/http"
)

type BaseApi struct {
	Ctx     *gin.Context
	Logger  *zap.Logger
	aborted bool
}

func NewBaseApi() *BaseApi {
	return &BaseApi{
		Logger: global.Logger,
	}
}

type BindDataOption struct {
	DTO               any
	BindParamsFromUri bool
}

func Build(c *gin.Context) *BaseApi {
	return NewBaseApi().BuildRequest(c)
}

func (m *BaseApi) BuildRequest(c *gin.Context) *BaseApi {
	//綁定上下文
	m.Ctx = c
	return m
}

func (m *BaseApi) Bind(option BindDataOption) *BaseApi {
	if option.DTO == nil {
		return m
	}

	// 綁定數據 (僅在 DTO 不為 nil 時執行)
	var errBinding error
	if option.BindParamsFromUri {
		errBinding = m.Ctx.ShouldBindUri(option.DTO)
	} else {
		errBinding = m.Ctx.ShouldBind(option.DTO)
	}

	if errBinding != nil {
		m.handleBindingError(errBinding)
	}

	return m
}

func (m *BaseApi) handleBindingError(errBinding error) {
	m.aborted = true
	response := ResponseJson{Status: http.StatusBadRequest}

	// JSON 型態錯誤
	var typeErr *json.UnmarshalTypeError
	if errors.As(errBinding, &typeErr) {
		global.Logger.Info(consts.ErrBindingType, zap.String("request_id", m.RequestID()), zap.Error(errBinding))
		response.Code = consts.ErrBindingType
		response.Msg = consts.MsgBindingTypeError
		response.Data = gin.H{"field": typeErr.Field, "expected": typeErr.Type.String()}
		m.Fail(response)
		return
	}

	// 驗證錯誤
	var validationErr validator.ValidationErrors
	if errors.As(errBinding, &validationErr) {
		global.Logger.Info(consts.ErrBindingValidate, zap.String("request_id", m.RequestID()), zap.Error(errBinding))
		errorsDetail := make([]gin.H, 0)
		for _, ve := range validationErr {
			errorsDetail = append(errorsDetail, gin.H{"field": ve.Field(), "rule": ve.Tag()})
		}
		response.Code = consts.ErrBindingValidate
		response.Msg = consts.MsgBindingValidateError
		response.Data = errorsDetail
		m.Fail(response)
		return
	}

	// 其他錯誤
	global.Logger.Warn(consts.ErrBinding, zap.String("request_id", m.RequestID()), zap.Error(errBinding))
	response.Code = consts.ErrBinding
	response.Msg = consts.MsgBindingError
	m.Fail(response)
	return
}

func (m *BaseApi) GetUserId() (IsLogin bool, userId uuid.UUID) {
	isLogin, ok := m.Ctx.Get("isLogin")
	if !ok || isLogin != true {
		return false, uuid.Nil
	}
	userIdStr, ok := m.Ctx.Get("userId")
	if !ok {
		return false, uuid.Nil
	}
	userId, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		return false, uuid.Nil
	}
	return true, userId
}

func (m *BaseApi) HasError() bool {
	return m.aborted
}

func (m *BaseApi) Context() context.Context {
	return m.Ctx.Request.Context()
}

func (m *BaseApi) RequestID() string {
	return m.Ctx.GetString("request_id")
}

func (m *BaseApi) OK(resp ResponseJson) {
	OK(m.Ctx, resp)
}

func (m *BaseApi) Fail(resp ResponseJson) {
	Fail(m.Ctx, resp)
}

func (m *BaseApi) ServerFail(resp ResponseJson) {
	ServerFail(m.Ctx, resp)
}
