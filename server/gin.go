package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Gin *gin.Engine
)

func init() {
	Gin = gin.Default()
	Gin.Use(CORSMiddleware())
}

func Run() {
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
}
