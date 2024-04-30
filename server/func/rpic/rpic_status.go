package rpic

import (
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"rpics-docker/server"
)

func rpicStatus(c *gin.Context) {
	stat := server.SQL.RpicStat.Get()
	totalImages := server.SQL.CountAllPics()
	c.JSON(
		http.StatusOK, gin.H{
			"code": 0,
			"msg":  "server running",
			"stat": gin.H{
				"total_images": totalImages,
				"cache_images": stat.CacheCount,
				"cache_space":  humanize.Bytes(uint64(stat.CacheSpace)),
				"image_space":  humanize.Bytes(uint64(stat.ImageSpace)),
			},
		},
	)
	return
}
