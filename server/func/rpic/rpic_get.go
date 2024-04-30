package rpic

import (
	"net/http"
	"net/url"
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
	var queries Queries
	var err error
	var queryValues url.Values
	var album string
	// parse queries
	if index := strings.LastIndex(context.Request.URL.Path, "!"); index != -1 {
		queryValues, err = url.ParseQuery(context.Request.URL.Path[index+1:])
		if err != nil {
			shortcut.RespStatusJSON(context, 1, "invalid query string")
			return
		}
		queries = Queries(queryValues)
		rawAlbum := context.Param("album")
		i := strings.LastIndex(rawAlbum, "!")
		if i > 0 {
			album = strings.Trim(rawAlbum[:i], "/")
		}
	} else {
		queries = Queries(context.Request.URL.Query())
		album = strings.Trim(context.Param("album"), "/")
	}
	// check queries
	if !checkQuery(context, queries) {
		return
	}
	// parse rid
	rawRid, ok := queries.Key("rid")
	rid, _ := strconv.Atoi(rawRid)
	rpicReq := &db.RpicRequest{
		Album:  album,
		Scale:  queries.Get("scale"),
		HasRid: ok,
		Rid:    rid,
	}
	quality, ok := image.Quality[queries.Get("quality")]
	if !ok {
		quality = 3
	}
	reqData := defaultData(
		&db.ImageData{
			Size:    queries.Get("size"),
			Quality: quality,
			Format:  queries.Get("format"),
		},
	)
	var main string
	reqData.Hash, main, ok = server.SQL.Rpic(rpicReq)
	if !ok {
		shortcut.RespStatusJSON(context, 1, "image not found")
		return
	}
	if _, ok = queries.Key("imageAve"); ok {
		context.JSON(
			http.StatusOK,
			gin.H{
				"code": 0, "main_color": str.Join("#", main), "id": reqData.Hash,
			},
		)
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
		contentSize = server.SQL.AddImageData(reqData)
		server.SQL.StatAddImageCache(contentSize, true)
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

type Queries url.Values

func (q Queries) Key(key string) (string, bool) {
	v, ok := q[key]
	if !ok {
		return "", false
	}
	return v[0], true
}

func (q Queries) Get(key string) string {
	v, ok := q.Key(key)
	if !ok {
		return ""
	}
	return v
}

func checkQuery(context *gin.Context, queries Queries) bool {
	if scale, ok := queries.Key("scale"); ok {
		if scale == "" {
			shortcut.RespStatusJSON(context, 1, str.T("invalid scale value: {scale}", scale))
			return false
		}
		scale := db.ParseScale(scale)
		if scale == "" {
			shortcut.RespStatusJSON(context, 1, str.T("invalid scale value: {scale}", scale))
			return false
		}
	}
	if rawRid, ok := queries.Key("rid"); ok {
		_, err := strconv.Atoi(rawRid)
		if err != nil {
			shortcut.RespStatusJSON(context, 1, str.T("invalid rid value: {rid}", rawRid))
			return false
		}
	}
	if size, ok := queries.Key("size"); ok {
		_, ok = image.Sizes[size]
		if !ok {
			shortcut.RespStatusJSON(context, 1, str.T("invalid size value: {size}", size))
			return false
		}
	}
	if quality, ok := queries.Key("quality"); ok {
		_, ok = image.Quality[quality]
		if !ok {
			shortcut.RespStatusJSON(context, 1, str.T("invalid quality value: {quality}", quality))
			return false
		}
	}
	if format, ok := queries.Key("format"); ok {
		_, ok = image.Formats[format]
		if !ok {
			shortcut.RespStatusJSON(context, 1, str.T("invalid format value: {format}", format))
			return false
		}
	}
	return true
}
