package serverlog

import (
	"git.realxlfd.cc/RealXLFD/golib/cli/logger"
	"github.com/gin-gonic/gin"
)

var (
	Log = logger.New(
		logger.WithLevel(logger.LevelDebug),
	)
	LogLevelMap = map[string]logger.Level{
		"silent": logger.LevelSlient,
		"error":  logger.LevelError,
		"warn":   logger.LevelWarn,
		"info":   logger.LevelInfo,
		"debug":  logger.LevelDebug,
	}
	GinLogLevelMap = map[string]string{
		"silent": gin.ReleaseMode,
		"error":  gin.ReleaseMode,
		"warn":   gin.ReleaseMode,
		"info":   gin.DebugMode,
		"debug":  gin.DebugMode,
	}
)
