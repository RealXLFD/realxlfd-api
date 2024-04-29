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
	switch album {
	case "":
		ok := server.SQL.RemoveAll(hash)
		if !ok {
			msg := str.T("can not remove image with hash({}) in the database", hash)
			log.Error(msg)
			shortcut.RespStatusJSON(context, 1, msg)
			return
		}
		log.Debug(str.T("remove image({hash}) from database successfully"), hash)
		shortcut.RespStatusJSON(context, 0, str.T("remove image({}) successfully", hash))
		return
	default:
		count, ok := server.SQL.CountAlbums(hash)
		if !ok {
			msg := "can not query database for albums contained image"
			log.Error(msg)
			shortcut.RespStatusJSON(context, 1, msg)
			return
		}
		switch count {
		case 2:
			// info: 从相册中移除图片
			server.SQL.RemoveFromAlbum(hash, album)
			shortcut.RespStatusJSON(
				context,
				0,
				str.T("remove image from album({}) successfully", album),
			)
			return
		case 1:
			// info: 删除所有图片
			var paths []string
			paths, ok = server.SQL.GetAllPaths(hash)
			if !ok {
				msg := str.T("can not get paths stored image({})", hash)
				log.Error(msg)
				shortcut.RespStatusJSON(context, 1, msg)
				return
			}
			for i := range paths {
				err := os.Remove(paths[i])
				if err != nil {
					msg := str.T("can not remove file: {path}", paths[i])
					log.Error(msg)
					shortcut.RespStatusJSON(context, 1, msg)
					return
				}
			}
			shortcut.RespStatusJSON(context, 0, "delete image successfully")
			return
		}
	}
	panic("fatal database error")
}
