package main

import (
	"os"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/cli/logger"
	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

func readConfig() (port string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("server.port", 80)
	viper.SetDefault("server.loglevel", "info")
	viper.SetDefault("rpic.concurrency", "auto")
	viper.SetDefault("rpic.throttling.limit", 100)
	viper.SetDefault("rpic.throttling.resetCycle", "1m")
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("please edit config.yaml to set up the server")
		err = viper.WriteConfigAs("config.yaml")
		if err != nil {
			log.Error("can not create config.yaml")
		}
		os.Exit(0)
	}
	// info: 设置日志级别
	loglevel := viper.GetString("server.loglevel")
	level, ok := serverlog.LogLevelMap[loglevel]
	if !ok {
		log.Warn("invalid log level: {level}", loglevel)
		log.Warn("use default log level: info")
		level = logger.LevelInfo
	} else {
		ginLevel := serverlog.GinLogLevelMap[loglevel]
		gin.SetMode(ginLevel)
	}
	serverlog.Log.Level = level
	// info: 读取端口
	rawPort := viper.GetInt("server.port")
	if rawPort == 0 {
		log.Error(str.T("invalid port: {port}", rawPort))
		os.Exit(1)
	}
	return str.Join(":", strconv.Itoa(rawPort))
}
