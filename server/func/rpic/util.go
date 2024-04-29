package rpic

import (
	"strconv"
	"time"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

func GetTimeStamp() (result int) {
	now := time.Now()
	timeStr := str.F("%s%03d", now.Format("20060102150405")[2:], now.Nanosecond()/1000000)
	result, _ = strconv.Atoi(timeStr)
	return
}
