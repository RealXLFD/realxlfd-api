package welcome

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rpics-docker/server/prehandler"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

func Serve(engine *gin.Engine) {
	engine.GET(
		"/", func(context *gin.Context) {
			context.JSON(
				http.StatusOK, gin.H{
					"msg": "welcome to api.realxlfd.cc",
					"available": []gin.H{
						{
							"/rpics": "get random image",
						},
					},
				},
			)
		},
	)
	// info: 身份验证测试
	engine.GET(
		"/auth", prehandler.Auth, func(context *gin.Context) {
			context.JSON(
				http.StatusOK, gin.H{
					"code": 0,
					"msg":  "authenticate successfully",
				},
			)
			return
		},
	)
	log.Info("welcome api loaded")
}
