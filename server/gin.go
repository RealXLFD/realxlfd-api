package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"rpics-docker/server/db"
	"rpics-docker/server/middleware"
	"rpics-docker/serverlog"
)

var (
	Gin  *gin.Engine
	SQL  *db.Sqlite
	ENV  = map[string]string{}
	log  = serverlog.Log
	Root = "./"
)

func Run(port string, apps ...func(engine *gin.Engine)) {
	Gin = gin.Default()
	Gin.Use(middleware.CORSMiddleware())
	token := os.Getenv("TOKEN")
	if token == "" {
		token = "realxlfd"
	}
	ENV["TOKEN"] = token
	SQL = db.Connect()
	for _, app := range apps {
		app(Gin)
	}
	// info: 从配置文件读取端口
	err := Gin.Run(port)
	if err != nil {
		log.Error("server failure: {err}", err.Error())
		os.Exit(1)
	}
}
