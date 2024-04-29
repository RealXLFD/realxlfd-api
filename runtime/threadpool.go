package server_runtime

import (
	"runtime"

	"git.realxlfd.cc/RealXLFD/golib/proc/thread/pool"
)

var (
	ThreadPool = pool.NewClosure(runtime.NumCPU(), runtime.NumCPU()).Run()
)
