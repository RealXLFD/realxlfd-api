package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"rpics-docker/server/db"
	"rpics-docker/server/image"
)

func rpicPostUpload(context *gin.Context) {
	ok := Auth(context)
	if !ok {
		return
	}
	album := strings.Trim(context.Param("album"), "/")
	form, err := context.MultipartForm()
	if err != nil {
		log.Error(
			"failed to get form data from request: {url}",
			context.Request.URL.Path,
		)
		context.JSON(http.StatusOK, gin.H{"code": "1", "msg": "invalid post request"})
		return
	}
	files := form.File["upload[]"]
	if len(files) == 0 {
		// err: 处理上传文件为空的情况
		context.JSON(http.StatusOK, gin.H{"code": "1", "msg": "no file uploaded"})
		return
	}
	uploadId := xid.New().String()
	uploads := make([]string, 0, len(files))
	defer func() {

	}()
	var finishes []gin.H
	skips, errs := &atomic.Int64{}, &atomic.Int64{}
	task := &sync.WaitGroup{}
	task.Add(len(files))
	for i, file := range files {
		path := filepath.Join(Root, "temp", str.Join(uploadId, "-", strconv.Itoa(i)))
		err = context.SaveUploadedFile(file, path)
		if err != nil {
			log.Error("failed to save uploaded file: {path}", path)
			errs.Add(1)
			task.Done()
			continue
		}
		log.Debug(str.T("receive and store temp file to {path}", path))
		uploads = append(uploads, path)
		ThreadPool.Push(
			func() {
				defer task.Done()
				hash := image.MD5(path)
				if SQL.Contains(hash) {
					skips.Add(1)
					return
				}
				pngDst := filepath.Join(
					Root,
					"/img",
					album,
					str.T("{hash}-raw-lossless.png", hash),
				)
				webpDst := filepath.Join(
					Root,
					"/img",
					album,
					str.T("{hash}-2k-Q=90.webp", hash),
				)
				ok := image.Convert(path, pngDst, "raw", 5)
				if !ok {
					errs.Add(1)
					return
				}
				meta := image.GetInfo(pngDst)
				ok = image.Convert(path, webpDst, "2k", 3)
				if !ok {
					errs.Add(1)
					err := os.Remove(pngDst)
					if err != nil {
						log.Error(str.T("failed to remove file: {dst}"), pngDst)
						return
					}
				}
				SQL.AddImage(
					&db.Image{
						Hash: hash, Main: meta.Prominent, Scale: meta.Scale,
						Date: GetTimeStamp(),
					},
				)
				SQL.AddImageData(
					&db.ImageData{
						Path: pngDst, Hash: hash, Size: "raw", Quality: 5, Format: "png",
					},
				)
				SQL.AddImageData(
					&db.ImageData{
						Path: webpDst, Hash: hash, Size: "2k", Quality: 3, Format: "webp",
					},
				)
				SQL.AddAlbum(hash, album)
				finishes = append(
					finishes, gin.H{
						"hash":       hash,
						"main_color": meta.Prominent,
						"width":      meta.Width,
						"height":     meta.Height,
					},
				)
			},
		)
	}
	task.Wait()
	if errs.Load() == 0 {
		context.JSON(
			http.StatusOK,
			gin.H{
				"code":   0,
				"msg":    str.T("successfully upload {count} images", len(finishes)),
				"images": finishes,
			},
		)
		return
	}
	context.JSON(
		http.StatusOK, gin.H{
			"code": 2,
			"msg": str.T(
				"finish processing {count} images, failed {count} images, skipped {count} images",
				len(finishes),
				errs.Load(),
				skips.Load(),
			),
			"images": finishes,
		},
	)
}
