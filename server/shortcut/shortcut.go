package shortcut

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespStatusJSON(context *gin.Context, code int, msg string) {
	context.JSON(
		http.StatusOK, gin.H{
			"code": code, "msg": msg,
		},
	)
}
