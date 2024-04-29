package server

import (
	"net/http"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
)

// auth info: 验证方式：http header auth标头，或url query的 key
func auth(context *gin.Context) {
	token, hasKey := context.GetQuery("key")
	if hasKey && token == ENV["TOKEN"] {
		return
	}
	token = context.GetHeader("auth")
	if token == ENV["TOKEN"] {
		return
	}
	context.AbortWithStatusJSON(
		http.StatusOK, gin.H{
			"code": 1,
			"msg":  str.T("authenticate failed: invalid token {token}", token),
		},
	)
}
