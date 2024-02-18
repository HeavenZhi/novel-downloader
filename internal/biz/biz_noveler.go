package biz

type BizNoveler interface {
	// 小说下载业务
	NovelDownload(getNovelMainUrl string, proxyList, userAgentList []string, saveNovelPath string, downloadSleepTime int64) (err error)
	// 小说主页面
	novelMain(getNovelMainUrl string, proxyList, userAgentList []string, saveNovelPath string) (err error)
	// 小说章节
	novelChapters(chapterURL string, proxyList, userAgentList []string, saveNovelPath, chapterTitle string) (err error)
}
