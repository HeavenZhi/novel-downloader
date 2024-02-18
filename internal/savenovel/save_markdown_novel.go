package savenovel

import (
	"bufio"
	"os"

	"github.com/HeavenZhi/novel-downloader/internal/config"
	"github.com/HeavenZhi/novel-downloader/internal/model"
	"go.uber.org/zap"
)

type SaveMarkdownNovel struct {
	style *config.SaveMarkdownStyleConfig
	// 写入器
	writer *bufio.Writer
}

// NewSaveNovel 创建小说保存器（构造函数）
func NewSaveMarkdownNovel() *SaveMarkdownNovel {
	return &SaveMarkdownNovel{
		style: config.Conf.SaveMarkdownStyleConfig,
	}
}

// CreateNovelFile 创建小说文件
func (s *SaveMarkdownNovel) CreateNovelFile(getNovelMainUrl, saveNovelPath string, novel *model.Novel) (err error) {
	// 获取小说文件路径
	file, err := os.OpenFile(saveNovelPath+novel.NovelTitle+".md", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		zap.L().Error("[Save Novel Chapters]", zap.Error(err))
		return err
	}

	//延迟执行，关闭文件
	defer file.Close()

	//接收一个 io.Writer ，返回包含默认缓存区大小为 4096 字节的 Writer 指针
	s.writer = bufio.NewWriter(file)

	// 写入小说标题
	s.writer.WriteString(s.style.NovelTitleLeftStyle + novel.NovelTitle + s.style.NovelTitleRightStyle + s.style.NewLineStyle)

	// 写入小说作者
	s.writer.WriteString(s.style.NewLineStyle + s.style.NovelAuthorStyle + novel.NovelAuthor + s.style.NewLineStyle)

	//写入小说简介样式
	s.writer.WriteString(s.style.NewLineStyle + s.style.NovelAbstractsStyle + s.style.NewLineStyle)
	for _, novelAbstracts := range novel.NovelAbstractsList {
		// 写入小说简介
		s.writer.WriteString(novelAbstracts + s.style.NewLineStyle)
	}

	// 写入小说链接
	s.writer.WriteString(s.style.NewLineStyle + s.style.NovelSoureStyle + getNovelMainUrl + s.style.NewLineStyle)

	// 写入我的版权信息
	s.writer.WriteString(s.style.NewLineStyle + "By: HeavenZhi's NovelDownload <https://github.com/HeavenZhi/novel-downloader>" + s.style.NewLineStyle)

	//将缓存数据刷新到底层的 io.Writer 对象
	defer s.writer.Flush()

	zap.L().Info("[Create Novel File]: Successfully create novel file: " + saveNovelPath + novel.NovelTitle + ".md")

	return nil
}

// SaveNovelChapters 保存小说章节
func (s *SaveMarkdownNovel) SaveNovelChapters(saveNovelPath, novelTitle, chapterTitle string, contentList []string) (err error) {
	// 获取小说文件路径
	file, err := os.OpenFile(saveNovelPath+novelTitle+".md", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		zap.L().Error("[Save Novel Chapters]", zap.Error(err))
		return err
	}

	// 延迟执行，关闭文件
	defer file.Close()
	// 接收一个 io.Writer ，返回包含默认缓存区大小为 4096 字节的 Writer 指针
	s.writer = bufio.NewWriter(file)

	// 写入小说章节标题
	s.writer.WriteString(s.style.NewLineStyle + s.style.TitleSeparatorStyle + chapterTitle + s.style.NewLineStyle)
	// 循环录入
	for _, content := range contentList {
		// 写入小说章节内容
		s.writer.WriteString(s.style.NewLineStyle + content + s.style.NewLineStyle)
	}

	// 将缓存数据刷新到底层的 io.Writer 对象
	defer s.writer.Flush()

	zap.L().Info("[Save Novel Chapters]: Successful Preservation of novel chapters: " + chapterTitle + ".")

	return nil
}
