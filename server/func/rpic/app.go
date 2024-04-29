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
	engine.POST("/rpic/add/:album", prehandler.Auth, rpicPOSTUpload)
	engine.PUT("/rpic/add/:album", prehandler.Auth, rpicPUTUpload)
	engine.GET(
		"/rpic/delete/:hash", prehandler.Auth, rpicDelete,
	)
	engine.GET(
		"/rpic/delete/:hash/:album", prehandler.Auth, rpicDelete,
	)
	engine.GET(
		"/rpic/get/:album", throttling, checkQuery, reqRpic,
	)
	engine.GET(
		"/rpic/get/:album/:rid", throttling, checkQuery, reqRpic,
	)
	engine.GET(
		"/rpic/get", throttling, checkQuery, reqRpic,
	)
	engine.GET("/rpic", throttling, checkQuery, reqRpic)
	log.Info("rpic api loaded")
}
