package handler

import (
	"github.com/HeavenZhi/novel-downloader/internal/biz"
)

// HandlerNovel 实现 HandlerNoveler 接口
type HandlerNovel struct {
	// 获取小说下载的业务接口对象
	novelDownloader biz.BizNoveler
}

// HandlerNovel 创建小说处理器对象(构造方法)
func NewHandlerNovel() *HandlerNovel {
	return &HandlerNovel{
		novelDownloader: biz.NewBizNovel(), // 创建小说下载业务对象
	}
}

// HandlerNovelInfo 处理小说信息
func (h *HandlerNovel) HandlerNovelInfo(getNovelMainUrl string, proxyList, userAgentList []string, saveNovelPath string, downloadSleepTime int64) (err error) {
	// 检查小说的主页地址是否为空
	if getNovelMainUrl == "" {
		return nil
	}

	// 小说下载业务
	err = h.novelDownloader.NovelDownload(getNovelMainUrl, proxyList, userAgentList, saveNovelPath, downloadSleepTime)
	if err != nil {
		return err
	}

	return nil
}
