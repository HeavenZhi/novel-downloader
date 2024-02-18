package savenovel

import (
	"github.com/HeavenZhi/novel-downloader/internal/model"
)

// SaveNoveler 保存小说
type SaveNoveler interface {
	// 创建小说文件
	CreateNovelFile(getNovelMainUrl, saveNovelPath string, novel *model.Novel) (err error)
	// 保存小说章节
	SaveNovelChapters(saveNovelPath, novelTitle, chapterTitle string, contentList []string) (err error)
}
