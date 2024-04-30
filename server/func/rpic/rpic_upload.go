package rpic

import (
	"io"
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
	"rpics-docker/server"
	"rpics-docker/server/db"
	"rpics-docker/server/image"
)

func rpicPOSTUpload(context *gin.Context) {
	album := strings.Trim(context.Param("album"), "/")
	form, err := context.MultipartForm()
	if err != nil {
		log.Error(
			"failed to get form data from request: {url}",
			context.Request.URL.Path,
		)
		context.JSON(http.StatusOK, gin.H{"code": 1, "msg": "invalid post request"})
		return
	}
	files := form.File["upload[]"]
	if len(files) == 0 {
		// err: 处理上传文件为空的情况
		context.JSON(http.StatusOK, gin.H{"code": 1, "msg": "no file uploaded"})
		return
	}
	uploadId := xid.New().String()
	uploads := make([]string, 0, len(files))
	var finishes []gin.H
	skips, errs := &atomic.Int64{}, &atomic.Int64{}
	task := &sync.WaitGroup{}
	task.Add(len(files))
	for i, file := range files {
		path := filepath.Join(
			CachePath,
			str.Join(uploadId, "-", strconv.Itoa(i)),
		)
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
				if server.SQL.Contains(hash, "") {
					if server.SQL.Contains(hash, album) {
						skips.Add(1)
						return
					}
					server.SQL.AddAlbum(hash, album)
					skips.Add(1)
					return
				}
				pngDst := storePath(hash, album, "raw", "lossless", "png")
				webpDst := storePath(hash, album, "2k", "Q=90", "webp")
				err = os.MkdirAll(filepath.Dir(pngDst), os.ModePerm)
				if err != nil {
					log.Error("can not create dir ./temp")
					errs.Add(1)
					return
				}
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
					}
					return
				}
				server.SQL.AddImage(
					&db.Image{
						Hash: hash, Main: meta.Prominent, Scale: meta.Scale,
						Date: GetTimeStamp(),
					},
				)
				server.SQL.AddImageData(
					&db.ImageData{
						Path: pngDst, Hash: hash, Size: "raw", Quality: 5, Format: "png",
					},
				)
				server.SQL.AddImageData(
					&db.ImageData{
						Path: webpDst, Hash: hash, Size: "2k", Quality: 3, Format: "webp",
					},
				)
				server.SQL.AddAlbum(hash, album)
				finishes = append(
					finishes, gin.H{
						"id":         hash,
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
	for i := range uploads {
		err := os.Remove(uploads[i])
		if err != nil {
			log.Error(str.T("failed to remove temp file: {path}"), uploads[i])
		}
	}
	context.JSON(
		http.StatusOK, gin.H{
			"code": 1,
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

func rpicPUTUpload(context *gin.Context) {
	album := strings.Trim(context.Param("album"), "/")
	path := filepath.Join(CachePath, xid.New().String())
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		msg := str.T("failed to write data to file: {path}", path)
		log.Error(msg)
		context.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg})
		return
	}
	defer func(path string) {
		err = os.Remove(path)
		if err != nil {
			log.Error(str.T("can not remove temp file: {}", path))
		}
	}(path)
	_, err = io.Copy(file, context.Request.Body)
	if err != nil {
		msg := str.T("failed to save file: {path}", path)
		log.Error(msg)
		context.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg})
		return
	}
	err = file.Close()
	if err != nil {
		log.Error(str.T("failed to close file: {path}", path))
	}
	log.Debug(str.T("receive and save temp file to {path}", path))
	hash := image.MD5(path)
	if server.SQL.Contains(hash, "") {
		if server.SQL.Contains(hash, album) {
			context.JSON(
				http.StatusOK,
				gin.H{"code": 2, "msg": str.T("the image already exist in {album}", album)},
			)
			return
		}
		server.SQL.AddAlbum(hash, album)
		context.JSON(
			http.StatusOK,
			gin.H{
				"code": 1,
				"msg":  str.T("add the image to album({}),but the image already exists", album),
			},
		)
		return
	}
	pngDst := storePath(hash, album, "raw", "lossless", "png")
	webpDst := storePath(hash, album, "2k", "Q=90", "webp")
	err = os.MkdirAll(filepath.Dir(pngDst), os.ModePerm)
	if err != nil {
		log.Error("can not create dir: {path}", filepath.Dir(pngDst))
		return
	}
	var ok bool
	task := sync.WaitGroup{}
	task.Add(1)
	ThreadPool.Push(
		func() {
			defer task.Done()
			ok = image.Convert(path, pngDst, "raw", 5)
		},
	)
	task.Wait()
	if !ok {
		context.JSON(
			http.StatusOK,
			gin.H{"code": 1, "msg": "invalid image: can not convert the image to PNG format"},
		)
		return
	}
	var meta image.Metadata
	task.Add(1)
	ThreadPool.Push(
		func() {
			defer task.Done()
			meta = image.GetInfo(pngDst)
			ok = image.Convert(path, webpDst, "2k", 3)
		},
	)
	task.Wait()
	if !ok {
		context.JSON(
			http.StatusOK,
			gin.H{"code": 1, "msg": "invalid image: can not convert the image to PNG format"},
		)
		err := os.Remove(pngDst)
		if err != nil {
			log.Error(str.T("failed to remove file: {dst}"), pngDst)
		}
		return
	}
	server.SQL.AddImage(
		&db.Image{
			Hash: hash, Main: meta.Prominent, Scale: meta.Scale,
			Date: GetTimeStamp(),
		},
	)
	contentSize := server.SQL.AddImageData(
		&db.ImageData{
			Path: pngDst, Hash: hash, Size: "raw", Quality: 5, Format: "png",
		},
	)
	server.SQL.StatAddImageCache(contentSize, false)
	contentSize = server.SQL.AddImageData(
		&db.ImageData{
			Path: webpDst, Hash: hash, Size: "2k", Quality: 3, Format: "webp",
		},
	)
	server.SQL.StatAddImageCache(contentSize, false)
	server.SQL.AddAlbum(hash, album)
	context.JSON(
		http.StatusOK, gin.H{
			"code": 0,
			"msg":  str.T("add image to {album} successfully", album),
			"images": []gin.H{
				{
					"id":         hash,
					"main_color": meta.Prominent,
					"width":      meta.Width,
					"height":     meta.Height,
				},
			},
		},
	)
	return
}

func storePath(hash, album, size, quality, format string) string {
	filename := str.T("{hash}-{size}-{quality}.{format}", hash, size, quality, format)
	return filepath.Join(ImagePath, album, filename)
}
