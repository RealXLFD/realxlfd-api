package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RpicServer() {
	Gin.POST(
		"/rpic/add/:album", func(context *gin.Context) {
			ok := Auth(context)
			if !ok {
				return
			}
			album := strings.Trim(context.Param("album"), "/")
			form, err := context.MultipartForm()
			if err != nil {
				log.Error(
					"failed to get form data from request: {url}",
					context.Request.URL.Path,
				)
			}
			files := form.File["files"]
			if len(files) == 0 {
				context.JSON(http.StatusOK, gin.H{"code": "1", "msg": "no file uploaded"})
				return
			}
			for i := range files {

			}

			// TODO: 处理图片上传
		},
	)
	Gin.GET(
		"/rpic/delete/:hash", func(context *gin.Context) {
			// TODO: 处理图片删除
		},
	)
	Gin.GET(
		"/rpic/delete/:hash/:album", func(context *gin.Context) {
			// TODO: 处理删除指定相册的图片
		},
	)
	Gin.GET(
		"/rpic/get/:album", func(context *gin.Context) {
			// TODO: rpics with album
		},
	)
	Gin.GET(
		"/rpic/get", func(context *gin.Context) {
			// TODO: rpics
		},
	)
}
