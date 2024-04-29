package prehandler

import (
	"net/http"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"rpics-docker/server"
)

// Auth info: 验证方式：http header auth标头，或url query的 key
func Auth(context *gin.Context) {
	token, hasKey := context.GetQuery("key")
	if hasKey && token == server.ENV["TOKEN"] {
		return
	}
	token = context.GetHeader("Auth")
	if token == server.ENV["TOKEN"] {
		return
	}
	context.AbortWithStatusJSON(
		http.StatusOK, gin.H{
			"code": 1,
			"msg":  str.T("authenticate failed: invalid token {token}", token),
		},
	)
}
