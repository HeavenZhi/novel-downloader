package biz

import (
	"time"

	"github.com/HeavenZhi/novel-downloader/internal/getdata"
	"github.com/HeavenZhi/novel-downloader/internal/model"
	"github.com/HeavenZhi/novel-downloader/internal/parsedata"
	"github.com/HeavenZhi/novel-downloader/internal/savenovel"
	"go.uber.org/zap"
)

// BizNovel 实现了 BizNoveler 接口
type BizNovel struct {
	// novel 小说模型
	novel *model.Novel
	// getNovel 获取小说数据
	getNovel getdata.GetNoveler
	// parseNovel 解析小说主页面数据
	parseNovel parsedata.ParseNoveler
	// saveNovel 保存小说数据
	saveNovel savenovel.SaveNoveler
}

// BizNovel 创建小说下载器（构造函数）
func NewBizNovel() *BizNovel {
	return &BizNovel{
		novel:      new(model.Novel),                  // 小说模型
		getNovel:   getdata.NewGetNovel(),             // 获取小说数据
		parseNovel: parsedata.NewParseQianyegeNovel(), // 解析 www.qianyege.com 小说主页面数据，可根据需要自行基于正则表达式实现 parsedata.ParseNoveler 接口，即可解析不同网站的小说数据
		saveNovel:  savenovel.NewSaveTxtNovel(),       // 保存小说数据，可根据需要自行实现 savenovel.SaveNovel 接口，即可保存指定格式的小说数据
	}
}

// NovelDownload 小说下载业务
func (b *BizNovel) NovelDownload(getNovelMainUrl string, proxyList, userAgentList []string, saveNovelPath string, downloadSleepTime int64) (err error) {
	// 获取小说主页面数据
	err = b.novelMain(getNovelMainUrl, proxyList, userAgentList, saveNovelPath)
	if err != nil {
		return err
	}

	// 获取小说的总章节数
	chaptersNumber := len(b.novel.NovelGeneralChaptersList)

	// 遍历小说章节数据
	for i := 0; i < chaptersNumber; i++ {

		// 休眠当前 goroutine，防止被封 IP
		time.Sleep(time.Duration(downloadSleepTime) * time.Second)

		// 获取小说章节数据
		err = b.novelChapters(b.novel.NovelGeneralChaptersList[i].ChapterURL, proxyList, userAgentList, saveNovelPath, b.novel.NovelGeneralChaptersList[i].ChapterTitle)
		if err != nil {
			return err
		}
	}

	return nil
}

// novelMain 小说主页面数据
func (b *BizNovel) novelMain(getNovelMainUrl string, proxyList, userAgentList []string, saveNovelPath string) (err error) {
	// 获取小说主页面数据
	rawData, err := b.getNovel.GetNovelData(getNovelMainUrl, proxyList, userAgentList)
	if err != nil {
		return err
	}

	// 解析小说主页面数据
	err = b.parseNovel.ParseNovelMain(rawData, b.novel, getNovelMainUrl)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		return err
	}

	// 保存小说主页数据
	err = b.saveNovel.CreateNovelFile(getNovelMainUrl, saveNovelPath, b.novel)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		return err
	}

	zap.L().Info("[Novel Main Data]: Novel main data done.")

	return nil
}

// novelChapters 小说章节数据
func (b *BizNovel) novelChapters(chapterURL string, proxyList, userAgentList []string, saveNovelPath, chapterTitle string) (err error) {
	// 获取小说章节数据
	rawData, err := b.getNovel.GetNovelData(chapterURL, proxyList, userAgentList)
	if err != nil {
		return err
	}

	// 解析小说章节数据
	contentList, err := b.parseNovel.ParseNovelChapters(rawData)
	if err != nil {
		return err
	}

	// 保存小说章节数据
	err = b.saveNovel.SaveNovelChapters(saveNovelPath, b.novel.NovelTitle, chapterTitle, contentList)
	if err != nil {
		return err
	}

	return nil
}
