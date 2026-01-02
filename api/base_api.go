package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/utils"
	"go.uber.org/zap"
)

const (
	ErrBindingType     = "BINDING_TYPE_ERROR"
	ErrBindingValidate = "BINDING_VALIDATE_ERROR"
	ErrBinding         = "BINDING_ERROR"
)
const (
	MsgBindingTypeError     = "binding.type_error"
	MsgBindingValidateError = "binding.Validate_error"
	MsgBindingError         = "Binding.error"
)

type BaseApi struct {
	Ctx    *gin.Context
	Errors error
	Logger *zap.Logger
}

func NewBaseApi() BaseApi {
	return BaseApi{
		Logger: global.Logger,
	}
}

type BuildRequestOption struct {
	Ctx               *gin.Context
	DTO               any
	BindParamsFromUri bool
}

func (m *BaseApi) BuildRequest(option BuildRequestOption) *BaseApi {
	var errBinding error
	/*===========================================================================================*/
	//綁定請求上下文
	m.Ctx = option.Ctx

	/*===========================================================================================*/
	//綁定請求的數據
	if option.DTO != nil {
		if option.BindParamsFromUri {
			errBinding = m.Ctx.ShouldBindUri(option.DTO)
		} else {
			errBinding = m.Ctx.ShouldBind(option.DTO)
		}
		/*===========================================================================================*/
		//如果嘗試綁定資料時發生錯誤
		if errBinding != nil {
			m.SetError(errBinding)
			//先創建一個要回傳的物件 錯誤代碼用400
			response := ResponseJson{Status: 400}

			//如果是Json轉型錯誤
			var typeErr *json.UnmarshalTypeError
			if errors.As(errBinding, &typeErr) {
				m.Logger.Info(ErrBindingType, zap.String("request_id", m.GetRequestID()), zap.Error(errBinding))
				response.Code = ErrBindingType
				response.Msg = MsgBindingTypeError
				response.Data = gin.H{
					"field":    typeErr.Field,
					"expected": typeErr.Type.String(),
					"value":    typeErr.Value,
				}
				return m.Fail(response)
			}

			//如果是資料驗證錯誤
			var validationErr validator.ValidationErrors
			if errors.As(errBinding, &validationErr) {
				m.Logger.Info(ErrBindingValidate, zap.String("request_id", m.GetRequestID()), zap.Error(errBinding))
				allValidationErrors := make([]gin.H, 0)
				for _, ve := range validationErr {
					allValidationErrors = append(allValidationErrors, gin.H{
						"field": ve.Field(),
						"rule":  ve.Tag(),
						"param": ve.Param(),
					})
				}
				response.Code = ErrBindingValidate
				response.Msg = MsgBindingValidateError
				response.Data = allValidationErrors
				return m.Fail(response)
			}

			//其他錯誤
			m.Logger.Warn(ErrBinding, zap.String("request_id", m.GetRequestID()), zap.Error(errBinding))
			response.Code = ErrBinding
			response.Msg = MsgBindingError

			return m.Fail(response)
		}

	}
	return m
}

func (m *BaseApi) SetError(errNew error) {
	m.Errors = utils.AppendError(m.Errors, errNew)
}

func (m *BaseApi) GetError() error {
	return m.Errors
}

func (m *BaseApi) GetRequestID() string {
	return m.Ctx.GetString("request_id")
}

//func (m *BaseApi) ParseValidateErrors(errs validator.ValidationErrors, target any) error {
//	var errResult error
//
//	var errValidation validator.ValidationErrors
//	ok := errors.As(errs, &errValidation)
//	//如果不是驗證的錯誤   就直接把錯誤說明返回
//	if !ok {
//		return errs
//	}
//
//	//通過反射獲取指針指向元素指定類型對象
//	fields := reflect.TypeOf(target).Elem()
//
//	for _, fieldErr := range errValidation {
//		//驗證發生錯誤的對象
//		field, _ := fields.FieldByName(fieldErr.Field())
//		//自定義的驗證錯誤訊息標籤
//		errMessageTag := fmt.Sprintf("%s_err", fieldErr.Tag())
//		//取得對應的錯誤訊息
//		errMessage := field.Tag.Get(errMessageTag)
//
//		if errMessage == "" {
//			//如果拿不到自定義的錯誤訊息  嘗試拿取message的錯誤訊息標籤(統一的錯誤訊息)
//			errMessage = field.Tag.Get("message")
//		}
//
//		if errMessage == "" {
//			//如果還是取不到錯誤訊息 就產生一個預設的錯誤訊息格式
//			errMessage = fmt.Sprintf("%s: %s Error (Validate fail)", fieldErr.Field(), fieldErr.Tag())
//		}
//
//		errResult = utils.AppendError(errResult, errors.New(errMessage))
//
//	}
//	return errResult
//}

func (m *BaseApi) OK(resp ResponseJson) *BaseApi {
	OK(m.Ctx, resp)
	return m
}

func (m *BaseApi) Fail(resp ResponseJson) *BaseApi {
	Fail(m.Ctx, resp)
	return m
}

func (m *BaseApi) ServerFail(resp ResponseJson) *BaseApi {
	ServerFail(m.Ctx, resp)
	return m
}
