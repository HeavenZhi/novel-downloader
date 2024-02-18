package logger

import "github.com/HeavenZhi/novel-downloader/internal/config"

type Logger interface {
	Init(logCfg *config.LogConfig) (err error)
}
