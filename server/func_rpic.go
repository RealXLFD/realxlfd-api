package server

import (
	"github.com/gin-gonic/gin"
)

func RpicServer() {
	Gin.POST("/rpic/add/:album", rpicPostUpload)
	Gin.PUT("/rpic/add")
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
