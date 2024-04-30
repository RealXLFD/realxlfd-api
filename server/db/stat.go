package db

import (
	"sync"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

type RpicStat struct {
	lock         *sync.Mutex
	requestCount int
	cacheCount   int
	cacheSpace   int64
	imageSpace   int64
}

func (db *Sqlite) StatAddRequest() {
	stat := db.rpicStat
	stat.lock.Lock()
	defer stat.lock.Unlock()
	stat.requestCount++
	_, err := db.driver.Exec("UPDATE RpicStat SET Request = ?", stat.requestCount)
	if err != nil {
		log.Error(str.T("sql error: failed to update statistic table"))
		return
	}
	log.Debug("update RpicStat successfully")
	return
}

func (db *Sqlite) StatAddImageCache(contentSize int64, isCache bool) {
	stat := db.rpicStat
	stat.lock.Lock()
	defer stat.lock.Unlock()
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
		return
	}
	log.Debug("update RpicStat successfully")
	return
}
