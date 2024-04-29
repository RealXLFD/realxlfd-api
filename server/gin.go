package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"rpics-docker/server/db"
	"rpics-docker/serverlog"
)

var (
	Gin  *gin.Engine
	SQL  = db.Connect()
	ENV  = map[string]string{}
	log  = serverlog.Log
	Root = "./"
)

func init() {
	Gin = gin.Default()
	Gin.Use(CORSMiddleware())
	token := os.Getenv("TOKEN")
	if token == "" {
		token = "realxlfd"
	}
	ENV["TOKEN"] = token
}

func ServeWelcome() {
	Gin.GET(
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
	Gin.GET(
		"/auth", auth, func(context *gin.Context) {
			context.JSON(
				http.StatusOK, gin.H{
					"code": 0,
					"msg":  "authenticate successfully",
				},
			)
			return
		},
	)
}

func Run() {
	err := Gin.Run(":80")
	if err != nil {
		log.Error("server failure: {err}", err.Error())
		os.Exit(1)
	}
}

func respStatusJSON(context *gin.Context, code int, msg string) {
	context.JSON(
		http.StatusOK, gin.H{
			"code": code, "msg": msg,
		},
	)
}
