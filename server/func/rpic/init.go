package rpic

import (
	"runtime"
	"strconv"
	"time"

	"git.realxlfd.cc/RealXLFD/golib/net/middleware/throttler"
	"git.realxlfd.cc/RealXLFD/golib/proc/thread/pool"
	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	ThreadPool *pool.ClosurePool
)

func Init() {
	// info: 读取配置文件中的并发数配置
	concurrency := viper.GetString("rpic.concurrency")
	switch concurrency {
	case "auto":
		ThreadPool = pool.NewClosure(runtime.NumCPU(), runtime.NumCPU()).Run()
	default:
		threads, err := strconv.Atoi(concurrency)
		if err != nil {
			panic(str.T("invalid concurrency: {concurrency}", concurrency))
		}
		ThreadPool = pool.NewClosure(threads, threads*2).Run()
	}
	// info: 读取配置文件中的限流配置
	rawResetCycle := viper.GetString("rpic.throttling.resetCycle")
	resetCycle, err := time.ParseDuration(rawResetCycle)
	if err != nil {
		panic(str.T("invalid reset cycle: {cycle}", rawResetCycle))
	}
	limit := viper.GetUint64("rpic.throttling.limit")
	Throttling(resetCycle, limit, resetCycle*5)
}

var (
	Throttler *throttler.Controller
)

func Throttling(resetCycle time.Duration, limit uint64, cleanCycle time.Duration) {
	Throttler = throttler.New(resetCycle, limit)
	Throttler.AutoClean(cleanCycle)
}

func throttling(context *gin.Context) {
	if Throttler == nil {
		return
	}
	if !Throttler.Add(context.ClientIP()) {
		context.AbortWithStatus(429)
		return
	}
}
