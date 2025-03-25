package logging

import "github.com/gin-gonic/gin"

const RequestIdKey = "logger_request_id"

func SetCtxRequestId(ginCtx *gin.Context, requestId string) {
	ginCtx.Set(RequestIdKey, requestId)
}
