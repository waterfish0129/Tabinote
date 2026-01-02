package api

import (
	"reflect"
)

type ResponseJson struct {
	Status    int    `json:"-"`
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data"`
	RequestID string `json:"requestId"`
}

func (m *ResponseJson) IsEmpty() bool {
	return reflect.DeepEqual(m, ResponseJson{})
}
