package rpic

import (
	"os"
	"strings"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"rpics-docker/server"
	"rpics-docker/server/shortcut"
)

func rpicDelete(context *gin.Context) {
	hash := strings.Trim(context.Param("hash"), "/")
	album := strings.Trim(context.Param("album"), "/")
	if !server.SQL.Contains(hash, "") {
		shortcut.RespStatusJSON(
			context,
			1,
			str.T("can not found image with hash({}) in the database", hash),
		)
		return
	}
	var count int
	var ok bool
	if album != "" {
		count, ok = server.SQL.CountAlbums(hash)
		if !ok {
			msg := "can not query database for albums contained image"
			log.Error(msg)
			shortcut.RespStatusJSON(context, 1, msg)
			return
		}
	}
	if count > 1 {
		// info: 从相册中移除图片
		server.SQL.RemoveFromAlbum(hash, album)
		shortcut.RespStatusJSON(
			context,
			0,
			str.T("remove image from album({}) successfully", album),
		)
		return
	} else if count == 1 || count == 0 {
		// info: 删除所有图片
		paths, ok := server.SQL.GetAllPaths(hash)
		if !ok {
			msg := str.T("can not get paths stored image({})", hash)
			log.Error(msg)
			shortcut.RespStatusJSON(context, 1, msg)
			return
		}
		for i := range paths {
			contentSize, err := os.Stat(paths[i])
			if err != nil {
				msg := str.T("can not get file stat: {path}", paths[i])
				log.Error(msg)
				shortcut.RespStatusJSON(context, 1, msg)
				return
			}
			err = os.Remove(paths[i])
			if err != nil {
				msg := str.T("can not remove file: {path}", paths[i])
				log.Error(msg)
				shortcut.RespStatusJSON(context, 1, msg)
				return
			}
			server.SQL.StatAddImageCache(-contentSize, strings.Contains())
		}
		shortcut.RespStatusJSON(context, 0, "delete image successfully")
		return
	}
	panic("fatal database error")
}
