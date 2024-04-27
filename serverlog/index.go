package serverlog

import "git.realxlfd.cc/RealXLFD/golib/cli/logger"

var (
	Log = logger.New(
		logger.WithLevel(logger.LevelDebug),
	)
)
