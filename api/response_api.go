package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func buildStatus(resp ResponseJson, defaultStatus int) int {
	if resp.Status == 0 {
		return defaultStatus
	}
	return resp.Status
}

func httpResponse(c *gin.Context, status int, resp ResponseJson) {
	if resp.IsEmpty() {
		c.AbortWithStatus(status)
		return
	}
	resp.RequestID = c.GetString("request_id")
	c.AbortWithStatusJSON(status, resp)
}

func OK(ctx *gin.Context, resp ResponseJson) {
	httpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}

func Fail(ctx *gin.Context, resp ResponseJson) {
	httpResponse(ctx, buildStatus(resp, http.StatusBadRequest), resp)
}

func ServerFail(ctx *gin.Context, resp ResponseJson) {
	httpResponse(ctx, buildStatus(resp, http.StatusInternalServerError), resp)
}
