package server

import (
	"os"
	"strings"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
)

func rpicDelete(context *gin.Context) {
	hash := strings.Trim(context.Param("hash"), "/")
	album := strings.Trim(context.Param("album"), "/")
	if !SQL.Contains(hash, "") {
		respStatusJSON(context, 1, str.T("can not found image with hash({}) in the database", hash))
		return
	}
	switch album {
	case "":
		ok := SQL.RemoveAll(hash)
		if !ok {
			msg := str.T("can not remove image with hash({}) in the database", hash)
			log.Error(msg)
			respStatusJSON(context, 1, msg)
			return
		}
		log.Debug(str.T("remove image({hash}) from database successfully"), hash)
		respStatusJSON(context, 0, str.T("remove image({}) successfully", hash))
		return
	default:
		count, ok := SQL.CountAlbums(hash)
		if !ok {
			msg := "can not query database for albums contained image"
			log.Error(msg)
			respStatusJSON(context, 1, msg)
			return
		}
		switch count {
		case 2:
			// info: 从相册中移除图片
			SQL.RemoveFromAlbum(hash, album)
			respStatusJSON(context, 0, str.T("remove image from album({}) successfully", album))
			return
		case 1:
			// info: 删除所有图片
			var paths []string
			paths, ok = SQL.GetAllPaths(hash)
			if !ok {
				msg := str.T("can not get paths stored image({})", hash)
				log.Error(msg)
				respStatusJSON(context, 1, msg)
				return
			}
			for i := range paths {
				err := os.Remove(paths[i])
				if err != nil {
					msg := str.T("can not remove file: {path}", paths[i])
					log.Error(msg)
					respStatusJSON(context, 1, msg)
					return
				}
			}
			respStatusJSON(context, 0, "delete image successfully")
			return
		}
	}
	panic("fatal database error")
}
