package handler

// HandlerNoveler 小说处理器接口
type HandlerNoveler interface {
	//处理小说信息
	HandlerNovelInfo(getNovelMainUrl string, proxyList, userAgentList []string, saveNovelPath string, downloadSleepTime int64) (err error)
}
