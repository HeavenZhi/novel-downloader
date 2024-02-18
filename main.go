package main

import (
	"flag"

	"github.com/HeavenZhi/novel-downloader/internal/config"
	"github.com/HeavenZhi/novel-downloader/internal/handler"
	"github.com/HeavenZhi/novel-downloader/internal/logger"
	"go.uber.org/zap"
)

func main() {
	//用于接收终端传入的参数
	var cfn string
	/*
		从命令行获取可能的conf路径（flag名, 默认值, 帮助信息）
			例如：novel-downloader -conf "./conf/config.toml"
	*/
	flag.StringVar(&cfn, "conf", "./conf/config.toml", "指定项目的配置文件路径")
	//解析命令行参数
	flag.Parse()

	// 接口的类型，实现类的对象，方便切换配置文件读取方式：创建配置文件对象
	var cfg config.Configer = &config.ViperConfig{}
	// 接口的类型，实现类的对象，方便切换配置文件读取方式：创建日志对象
	var log logger.Logger = &logger.ZapLogger{}

	// 0.加载配置文件
	err := cfg.Init(cfn)
	if err != nil {
		//程序启动时加载配置文件失败直接退出
		panic(err)
	}

	// 1.使用配置文件来初始日志模块
	err = log.Init(config.Conf.LogConfig)
	if err != nil {
		//程序启动时加载配置文件失败直接退出
		panic(err)
	}

	// 接口的类型，实现类的对象，方便切换配置文件读取方式：创建小说处理对象
	var hand handler.HandlerNoveler = handler.NewHandlerNovel()

	//处理小说信息
	err = hand.HandlerNovelInfo(config.Conf.NovelConfig.GetNovelMainUrl, config.Conf.SettingConfig.ProxyList, config.Conf.SettingConfig.UserAgentList, config.Conf.NovelConfig.SaveNovelPath, config.Conf.SettingConfig.DownloadSleepTime)
	if err != nil {
		//处理小说信息失败直接退出
		zap.L().Error("[Error Handle Novel Info]: ", zap.Error(err))
		return
	}
}
