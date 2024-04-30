package rpic

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"rpics-docker/server"
	"rpics-docker/server/prehandler"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

var (
	ImagePath = filepath.Join(server.Root, "/rpic/image")
	CachePath = filepath.Join(server.Root, "/rpic/temp")
)

func Serve(engine *gin.Engine) {
	Init()
	err := os.MkdirAll(CachePath, os.ModePerm)
	if err != nil {
		log.Error("can not create dir ./temp")
	}
	apiGroup := engine.Group("/rpic")
	{
		apiGroup.POST("/add/:album", prehandler.Auth, rpicPOSTUpload)
		apiGroup.PUT("/add/:album", prehandler.Auth, rpicPUTUpload)
		apiGroup.GET(
			"/delete/:hash", prehandler.Auth, rpicDelete,
		)
		apiGroup.GET(
			"/delete/:hash/:album", prehandler.Auth, rpicDelete,
		)
		apiGroup.GET(
			"/get/:album", throttling, rpicReq,
		)
		apiGroup.GET(
			"/get", throttling, rpicReq,
		)
		apiGroup.GET("/", throttling, rpicReq)
		apiGroup.GET("/status", rpicStatus)
		apiGroup.GET("/list/:album", rpicList)
		apiGroup.GET("/list", rpicList)
	}
	log.Info("rpic api loaded")
}
