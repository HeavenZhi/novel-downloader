package model

// Novel 小说模型对象
type Novel struct {
	NovelTitle               string
	NovelAuthor              string
	NovelAbstractsList       []string
	NovelGeneralChaptersList []*NovelGeneralChapter
}

// NovelGeneralChapter 小说章节模型对象
type NovelGeneralChapter struct {
	ChapterTitle string
	ChapterURL   string
}
