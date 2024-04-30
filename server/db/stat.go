package db

import (
	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

type RpicStat struct {
	requestCount int
	cacheCount   int
	cacheSpace   int64
	imageSpace   int64
	runner       bool
	RequestAdder chan struct{}
	CacheAdder   chan struct {
		contentSize int64
		isCache     bool
	}
}

func (stat *RpicStat) Get() struct {
	RequestCount, CacheCount int
	CacheSpace, ImageSpace   int64
} {
	return struct {
		RequestCount, CacheCount int
		CacheSpace, ImageSpace   int64
	}{
		RequestCount: stat.requestCount,
		CacheCount:   stat.cacheCount,
		CacheSpace:   stat.cacheSpace,
		ImageSpace:   stat.imageSpace,
	}
}

func (db *Sqlite) RunStat() {
	stat := db.RpicStat
	stat.runner = true
	go func() {
		for stat.runner {
			select {
			case <-stat.RequestAdder:
				stat.requestCount++
				_, err := db.driver.Exec("UPDATE RpicStat SET Request = ?;", stat.requestCount)
				if err != nil {
					log.Error(str.T("sql error: failed to update statistic table"))
				} else {
					log.Debug("update RpicStat successfully")
				}
			case cacheAdder := <-stat.CacheAdder:
				isCache, contentSize := cacheAdder.isCache, cacheAdder.contentSize
				if isCache {
					if contentSize > 0 {
						stat.cacheCount++
					} else {
						stat.cacheCount--
					}
					stat.cacheSpace += contentSize
				} else {
					stat.imageSpace += contentSize
				}
				_, err := db.driver.Exec(
					"UPDATE RpicStat SET CacheCount = ? , ImageSpace = ?,CacheSpace = ?;",
					stat.cacheCount,
					stat.imageSpace,
					stat.cacheSpace,
				)
				if err != nil {
					log.Error(str.T("sql error: failed to update statistic table"))
				} else {
					log.Debug("update RpicStat successfully")
				}
			}
		}
	}()
}

func (db *Sqlite) StatAddRequest() {
	db.RpicStat.RequestAdder <- struct{}{}
}

func (db *Sqlite) StatAddImageCache(contentSize int64, isCache bool) {
	db.RpicStat.CacheAdder <- struct {
		contentSize int64
		isCache     bool
	}{contentSize, isCache}
}
