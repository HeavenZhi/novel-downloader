/*
	Copyright 2024 HeavenZhi. All rights reserved.
	Use of this source code is governed by a MIT style
	license that can be found in the LICENSE file.
*/

package logger

import (
	"io"
	"os"

	"github.com/HeavenZhi/novel-downloader/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/*
 Zap 日志库三要素：
 	1.encoder 编码(日志格式、时间格式、字段格式)
 	2.日志的输出位置(终端、文件)
 	3.日志级别(需要输出的日志级别)
*/

type ZapLogger struct {
	// 声明 Logger 日志记录器
	logger *zap.Logger
}

// 获取普通编码器
func (z *ZapLogger) getEncoder() zapcore.Encoder {
	//zap.NewProductionEncoderConfig() 为生产环境返回一个固执己见的 EncoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()

	//zapcore.ISO8601TimeEncoder 时间编码器序列化时间，以毫秒精度到达 ISO8601 格式字符串的时间
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//zapcore.CapitalLevelEncoder 将 Level 序列化为全大写字符串，例如：InfoLevel 被序列化为“INFO”
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//zapcore.NewConsoleEncoder() 创建了一个编码器，其输出专为人类而非机器消费而设计，它以纯文本格式序列化核心日志条目数据（消息、级别、时间戳等），并将结构化上下文保留为 JSON
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 获取 lumberjack 日志分割编写器，设置日志记录路径
func (z *ZapLogger) getLogWriter(logCfg *config.LogConfig) zapcore.WriteSyncer {
	//lumberjack 日志分割
	lumberJackLogger := &lumberjack.Logger{
		//日志文件的保存位置
		Filename: logCfg.Filename,
		//日志文件在进行切割之前，日志文件的最大大小（以 MB 为单位）
		MaxSize: logCfg.MaxSize,
		//保留旧文件的最大个数
		MaxBackups: logCfg.MaxBackups,
		//保留旧文件的最大天数
		MaxAge: logCfg.MaxAge,
		//是否压缩/归档旧文件
		Compress: false,
	}

	//io.MultiWriter() 接收多个输入流，将这些输入流合并为一个输出流
	ws := io.MultiWriter(lumberJackLogger, os.Stdout)

	//zapcore.AddSync() 转换 io，写入到写入同步器。它试图变得智能：如果 io 的具体类型。Writer 实现了 WriteSyncer，我们将使用现有的 Sync 方法。如果没有，我们将添加一个无操作同步
	return zapcore.AddSync(ws)
}

// 获取 lumberjack 日志分割编写器，设置日志记录路径
func (z *ZapLogger) getErrorLogWriter(logCfg *config.LogConfig) zapcore.WriteSyncer {
	//lumberjack 日志分割
	lumberJackLogger := &lumberjack.Logger{
		//Error 日志文件的单独保存位置
		Filename: logCfg.ErrorFilename,
		//日志文件在进行切割之前，日志文件的最大大小（以 MB 为单位）
		MaxSize: logCfg.MaxSize,
		//保留旧文件的最大个数
		MaxBackups: logCfg.MaxBackups,
		//保留旧文件的最大天数
		MaxAge: logCfg.MaxAge,
		//是否压缩/归档旧文件
		Compress: false,
	}

	//io.MultiWriter() 接收多个输入流，将这些输入流合并为一个输出流
	ws := io.MultiWriter(lumberJackLogger)

	//zapcore.AddSync() 转换 io，写入到写入同步器。它试图变得智能：如果 io 的具体类型。Writer 实现了 WriteSyncer，我们将使用现有的 Sync 方法。如果没有，我们将添加一个无操作同步
	return zapcore.AddSync(ws)
}

// Init 初始化 Logger 日志记录器
func (z *ZapLogger) Init(logCfg *config.LogConfig) (err error) {
	//调用获取普通编码器
	encoder := z.getEncoder()

	//获取 lumberjack 日志分割编写器，设置日志记录路径
	writeSyncer := z.getLogWriter(logCfg)

	//获取 lumberjack 日志分割编写器，设置 Error 日志单独记录路径
	writeErrorSyncer := z.getErrorLogWriter(logCfg)

	//zapcore.NewCore() 创建一个包含普通编码器、日志编写器、默认日志级别为 Debug 的 WriteSyncer 核心
	debugCore := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//zapcore.NewCore() 创建一个包含普通编码器、日志编写器、默认日志级别为 Error 的 WriteSyncer 核心
	errCore := zapcore.NewCore(encoder, writeErrorSyncer, zap.ErrorLevel)

	/*
		zapcore.NewTee() 创建一个核心，用于将日志条目复制到两个或多个底层核心中
	*/
	core := zapcore.NewTee(debugCore, errCore)

	//zap.New() 从提供的 zapcore 构造一个新的 Logge 核心和选项。如果通过的 zapcore，Core 是 nil，它会退回到使用 no-op 实现
	//zap.AddCaller() 将记录器配置为使用 zap 调用方的文件名、行号和函数名称注释每条消息
	z.logger = zap.New(core, zap.AddCaller())

	//logger.Sugar() 包装了记录器，以提供更符合人体工程学但速度稍慢的 API，返回一个 SugaredLogger 指针
	//SugaredLogger 是非常便宜的，因此单个应用程序同时使用记录器和 SugaredLogger 是合理的，在性能敏感代码的边界上在它们之间进行转换
	// Logger = logger.Sugar()

	//替换 zap 全局 Logger
	zap.ReplaceGlobals(z.logger)
	//调用全局 Logger 示例
	zap.L().Info("[Zap Logger]: Zap logger init success")

	return
}
