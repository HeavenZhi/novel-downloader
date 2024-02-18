package parsedata

import "github.com/HeavenZhi/novel-downloader/internal/model"

// 解析小说
type ParseNoveler interface {
	// 解析小说主页面
	ParseNovelMain(data []byte, novel *model.Novel, getNovelMainUrl string) (err error)
	// 解析小说章节
	ParseNovelChapters(data []byte) (contentList []string, err error)
	// 解析源数据的编码
	parseSourceDataEncode(inData []byte) (outData []byte, err error)
	// 解析小说标题
	parseNovelTitle(dataString string) (novelTitle string, err error)
	// 解析小说作者
	parseNovelAuthor(dataString string) (novelAuthor string, err error)
	// 解析小说简介
	parseNovelAbstracts(dataString string) (novelAbstractsList []string, err error)
	// 解析小说总章节数
	parseNovelGeneralChapters(getNovelMainUrl, dataString string) (generalChapter []*model.NovelGeneralChapter, err error)
}
