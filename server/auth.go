package server

import (
	"net/http"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
)

// Auth info: 验证方式：http header auth标头，或url query的 key
func Auth(context *gin.Context) (ok bool) {
	token, hasKey := context.GetQuery("key")
	if hasKey && token == ENV["TOKEN"] {
		return true
	}
	token = context.GetHeader("auth")
	if token == ENV["TOKEN"] {
		return true
	}
	context.JSON(
		http.StatusOK, gin.H{
			"code": 1,
			"msg":  str.T("authenticate failed: invalid token {token}", token),
		},
	)
	return false
}
