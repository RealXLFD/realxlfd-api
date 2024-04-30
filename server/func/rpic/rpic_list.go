package rpic

import (
	"net/http"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"rpics-docker/server"
	"rpics-docker/server/shortcut"
)

func rpicList(c *gin.Context) {
	album := c.Param("album")
	rawPage := c.Query("page")
	if rawPage == "" {
		rawPage = "1"
	}
	page, err := strconv.Atoi(rawPage)
	if err != nil || page <= 0 {
		shortcut.RespStatusJSON(c, 1, str.T("invalid page number: {}", rawPage))
		return
	}
	if album == "" {
		shortcut.RespStatusJSON(c, 1, "album name is empty")
		return
	}

	count, ok := server.SQL.CountPicsByAlbum(album)
	if !ok {
		shortcut.RespStatusJSON(c, 1, "internal server error")
		return
	}
	if count == 0 {
		shortcut.RespStatusJSON(c, 1, str.T("empty album({})", album))
		return
	}
	ids := server.SQL.ListPics(album, 15, page)
	if ids == nil {
		shortcut.RespStatusJSON(c, 1, "internal server error")
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"code":         0,
			"msg":          str.T("find {} images in album({})", len(ids), album),
			"total_images": count,
			"images":       ids,
			"page":         page,
			"total_page":   (count + 14) / 15,
		},
	)
	return
}
