package server

import (
	"os"
)

func ServeRpic() {
	err := os.MkdirAll("./temp", os.ModePerm)
	if err != nil {
		log.Error("can not create dir ./temp")
	}
	Gin.POST("/rpic/add/:album", auth, rpicPOSTUpload)
	Gin.PUT("/rpic/add/:album", auth, rpicPUTUpload)
	Gin.GET(
		"/rpic/delete/:hash", auth, rpicDelete,
	)
	Gin.GET(
		"/rpic/delete/:hash/:album", auth, rpicDelete,
	)
	Gin.GET(
		"/rpic/get/:album", reqRpic,
	)
	Gin.GET(
		"/rpic/get", reqRpic,
	)
}
