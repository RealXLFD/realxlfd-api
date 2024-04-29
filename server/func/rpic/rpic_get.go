package rpic

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"rpics-docker/server"
	"rpics-docker/server/db"
	"rpics-docker/server/image"
	"rpics-docker/server/shortcut"
)

func reqRpic(context *gin.Context) {
	log.Info(
		str.T(
			"new rpic request({path}) from {ip}", context.Request.URL.Path, context.ClientIP(),
		),
	)
	album := strings.Trim(context.Param("album"), "/")
	rawRid, ok := context.GetQuery("rid")
	if !ok {
		rawRid = strings.Trim(context.Param("rid"), "/")
		if rawRid != "" {
			ok = true
		}
	}
	rid, _ := strconv.Atoi(rawRid)
	rpicReq := &db.RpicRequest{
		Album:  album,
		Scale:  context.Query("scale"),
		HasRid: ok,
		Rid:    rid,
	}
	quality, ok := image.Quality[context.Query("quality")]
	if !ok {
		quality = 3
	}
	reqData := defaultData(
		&db.ImageData{
			Size:    context.Query("size"),
			Quality: quality,
			Format:  context.Query("format"),
		},
	)
	var main string
	reqData.Hash, main, ok = server.SQL.Rpic(rpicReq)
	if !ok {
		shortcut.RespStatusJSON(context, 1, "image not found")
		return
	}
	var contentSize int64
	reqData.Path, contentSize, ok = server.SQL.GetPath(reqData)
	if !ok {
		shortcut.RespStatusJSON(context, 1, "internal server error: can not get image data")
		return
	}
	if contentSize == 0 {
		var pngPath string
		pngPath, _, ok = server.SQL.GetPath(
			&db.ImageData{
				Hash: reqData.Hash, Size: "raw", Quality: 5, Format: "png",
			},
		)
		if !ok {
			shortcut.RespStatusJSON(
				context,
				1,
				"internal server error: can not found original image",
			)
			return
		}
		task := sync.WaitGroup{}
		task.Add(1)
		ThreadPool.Push(
			func() {
				defer task.Done()
				reqData.Path = picConvert(pngPath, rpicReq.Album, reqData)
			},
		)
		task.Wait()
		reqData.Path = picConvert(pngPath, rpicReq.Album, reqData)
		server.SQL.AddImageData(reqData)
		stat, err := os.Stat(reqData.Path)
		if err != nil || stat.IsDir() {
			shortcut.RespStatusJSON(
				context,
				1,
				"internal server error: can not get content size of the converted image",
			)
			return
		}
		contentSize = stat.Size()
	}
	if _, ok = context.GetQuery("imageAve"); ok {
		context.JSON(
			http.StatusOK,
			gin.H{
				"code": 0, "main_color": str.Join("#", main), "id": reqData.Hash,
				"image": gin.H{"size": contentSize, "format": reqData.Format},
			},
		)
		return
	}
	file, err := os.Open(reqData.Path)
	if err != nil {
		shortcut.RespStatusJSON(
			context,
			1,
			"internal server error: can not get image data of the target image",
		)
		return
	}
	defer func(file *os.File, path string) {
		err := file.Close()
		if err != nil {
			log.Error(str.T("error occurred while closing the file: {}", path))
		}
	}(file, reqData.Path)
	log.Debug(str.T("handle rpic request with image({path})", reqData.Path))
	context.DataFromReader(
		http.StatusOK, contentSize, image.Formats[reqData.Format], file,
		map[string]string{
			"main-color": str.Join("#", main),
		},
	)
	return
}

// alert: 错误返回""
func picConvert(pngPath, album string, i *db.ImageData) (path string) {
	dst := storePath(i.Hash, album, i.Size, image.QualityArr[i.Quality], i.Format)
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		log.Error(str.T("can not create path: {path}", filepath.Dir(dst)))
	}
	ok := image.Convert(pngPath, dst, i.Size, i.Quality)
	if !ok {
		return ""
	}
	return dst
}

func defaultData(i *db.ImageData) *db.ImageData {
	if i.Size == "" {
		i.Size = "2k"
	}
	if i.Format == "" {
		i.Format = "webp"
	}
	return i
}

func checkQuery(context *gin.Context) {
	if scale, ok := context.GetQuery("scale"); ok {
		if scale == "" {
			shortcut.RespStatusJSON(context, 1, str.T("invalid scale value: {scale}", scale))
			context.Abort()
			return
		}
		scale := db.ParseScale(scale)
		if scale == "" {
			shortcut.RespStatusJSON(context, 1, str.T("invalid scale value: {scale}", scale))
			context.Abort()
		}
		return
	}
	if rawRid, ok := context.GetQuery("rid"); ok {
		_, err := strconv.Atoi(rawRid)
		if err != nil {
			shortcut.RespStatusJSON(context, 1, str.T("invalid rid value: {rid}", rawRid))
			context.Abort()
		}
		return
	}
	if size, ok := context.GetQuery("size"); ok {
		_, ok = image.Sizes[size]
		if !ok {
			shortcut.RespStatusJSON(context, 1, str.T("invalid size value: {size}", size))
			context.Abort()
		}
		return
	}
	if quality, ok := context.GetQuery("quality"); ok {
		_, ok = image.Quality[quality]
		if !ok {
			shortcut.RespStatusJSON(context, 1, str.T("invalid quality value: {quality}", quality))
			context.Abort()
		}
		return
	}
	if format, ok := context.GetQuery("format"); ok {
		_, ok = image.Formats[format]
		if !ok {
			shortcut.RespStatusJSON(context, 1, str.T("invalid format value: {format}", format))
			context.Abort()
		}
		return
	}
}
