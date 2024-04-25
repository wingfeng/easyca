package main

import (
	"easyca/conf"
	"easyca/engine"
	"easyca/rbac"

	//	_ "easyca/rbac"
	"easyca/routers"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AppName      string // 应用名称
	AppVersion   string // 应用版本
	BuildVersion string // 编译版本
	BuildTime    string // 编译时间
	GitRevision  string // Git版本
	GitBranch    string // Git分支
	GoVersion    string // Golang信息
)

func main() {
	showVersion := flag.Bool("ver", false, "程序版本")
	cfgFile := flag.String("c", "conf/config.yaml", "配置文件路径")
	flag.Parse()
	conf.AppName = AppName
	conf.AppVersion = AppVersion
	conf.BuildVersion = BuildVersion
	conf.BuildTime = BuildTime
	conf.GitRevision = GitRevision
	conf.GitBranch = GitBranch
	conf.GoVersion = GoVersion
	if *showVersion {
		conf.Version()
		return
	}

	//配置viper的配置文件路径
	config := conf.InitConfig(*cfgFile)
	logLevel := slog.LevelWarn
	switch strings.ToLower(conf.Default.Log) {
	case "debug":
		logLevel = slog.LevelDebug

	case "info":
		logLevel = slog.LevelInfo

	case "warn":
		logLevel = slog.LevelWarn

	}
	slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), logLevel)
	// //配置log输出级别
	// if !strings.EqualFold("", conf.Default.Log) {
	// 	//配置Log

	// 	logLevel, lex := log.LogLevelFromString(conf.Default.Log)
	// 	if !lex {
	// 		logLevel = log.DebugLvl
	// 	}
	// 	logger, _ := log.LoggerFromWriterWithMinLevel(os.Stdout, logLevel)
	// 	log.ReplaceLogger(logger)
	// }

	r := gin.Default()

	rbac.InitEnforcer()
	engine.InitEngine()

	router := routers.InitRouter(r)
	//Enable Cache

	router.Run(fmt.Sprintf("%s:%d", config.IP, config.Port))
}
