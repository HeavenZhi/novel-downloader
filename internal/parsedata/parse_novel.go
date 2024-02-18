package parsedata

import (
	"regexp"
	"strings"

	"github.com/HeavenZhi/novel-downloader/internal/encode"
	"github.com/HeavenZhi/novel-downloader/internal/model"
	"go.uber.org/zap"
)

// ParseQianyegeNovel 实现了 ParseNoveler 接口
type ParseQianyegeNovel struct {
	// 正则表达式对象
	re *regexp.Regexp
	// 编码器
	ecode encode.Encoder
}

// NewParseNovel 创建一个 ParseNoveler 对象（构造方法）
func NewParseQianyegeNovel() *ParseQianyegeNovel {
	return &ParseQianyegeNovel{
		ecode: encode.NewEncode(), // 编码器
	}
}

// ParseNovelMain 解析小说主页面
func (p *ParseQianyegeNovel) ParseNovelMain(data []byte, novel *model.Novel, getNovelMainUrl string) (err error) {
	// 获取源数据的编码方式
	data, err = p.parseSourceDataEncode(data)
	if err != nil {
		return err
	}

	// 将源数据转换为字符串
	dataString := string(data)

	// 解析小说标题
	novel.NovelTitle, err = p.parseNovelTitle(dataString)
	if err != nil {
		return err
	}

	// 解析小说作者
	novel.NovelAuthor, err = p.parseNovelAuthor(dataString)
	if err != nil {
		return err
	}

	// 解析小说简介
	novel.NovelAbstractsList, err = p.parseNovelAbstracts(dataString)
	if err != nil {
		return err
	}

	// 解析小说章节
	novel.NovelGeneralChaptersList, err = p.parseNovelGeneralChapters(getNovelMainUrl, dataString)
	if err != nil {
		return err
	}

	zap.L().Info("[Parse Novel Main]: Parse novel main is success.")

	return nil
}

// ParseNovelChapters 解析小说章节
func (p *ParseQianyegeNovel) ParseNovelChapters(data []byte) (contentList []string, err error) {
	// 获取源数据的编码方式
	data, err = p.parseSourceDataEncode(data)
	if err != nil {
		return nil, err
	}

	// 将源数据转换为字符串
	dataString := string(data)

	/*
		创建小说章节的正则表达式对象
			可以一次性获取到小说章节，且无需额外操作的正则表达式表达式：(?<=<div id="intro">\s*<p>&nbsp;&nbsp;&nbsp;&nbsp;).*(?=</p>\s*</div>)
			Go 不支持这个写法！！！
	*/
	p.re = regexp.MustCompile("(&nbsp;){4}[\\s\\S]*")

	// 获取小说章节的内容部分
	content := p.re.FindString(dataString)

	// 获取多余部分的开头
	discardIndex := strings.Index(content, "</div>")

	// 舍弃多余部分
	content = content[:discardIndex]

	// 去除无用的格式
	content = strings.ReplaceAll(content, "&nbsp;&nbsp;&nbsp;&nbsp;", "")

	// 必须替换掉换行符，这些换行符看不见，但是后面在设置输出格式时会形成巨大难度
	content = strings.ReplaceAll(content, "\r\n", "")

	// 按照 <br /><br /> 分割每一段文字
	contentList = strings.Split(content, "<br /><br />")

	zap.L().Info("[Parse Novel Chapters]: Parse novel chapters is success.")

	return contentList, nil
}

// ParseSourceDataEncode 获取源数据的编码
func (p *ParseQianyegeNovel) parseSourceDataEncode(inData []byte) (outData []byte, err error) {
	/*
		创建网站源代码编码的正则表达式对象
			可以一次性获取到字符编码，且无需额外操作的正则表达式表达式：(?<=(<meta.*charset="*))[^>"/]+
			Go 不支持这个写法！！！
	*/
	p.re = regexp.MustCompile("(<meta.*charset=\"*)[^\"/>]+")

	// 使用正则表达式匹配原始数据，获取到 meta 标签中的 charset 属性，以获取： <meta charset="字符编码">
	charsetString := p.re.FindString(string(inData))

	// 获取最后一个 = 符号的位置
	eqIndex := strings.LastIndex(charsetString, "=")

	// 获取字符编码
	charset := charsetString[eqIndex+1:]

	// 去除引号
	charset = strings.Trim(charset, "\"")
	// 去除空格
	charset = strings.Trim(charset, " ")

	// 判断是否需要对源数据进行转换编码
	if strings.ToUpper(charset) == "GBK" {
		// 将源数据转换为 GBK 编码
		outData, _ = p.ecode.EncodeToGBK(inData)
	} else {
		outData = inData
	}

	// 打印字符编码
	zap.L().Info("[Parse Source Data Encode]: The coding method of successfully parsing the novel data on the web site is: " + strings.ToUpper(charset) + ".")

	return outData, nil
}

// parseNovelTitle 解析小说标题
func (p *ParseQianyegeNovel) parseNovelTitle(dataString string) (novelTitle string, err error) {
	/*
		创建小说标题的正则表达式对象
			可以一次性获取到小说标题，且无需额外操作的正则表达式表达式：(?<=<h1.*>).*(?=</h1>)
			Go 不支持这个写法！！！
	*/
	p.re = regexp.MustCompile("<h1[\"=\\w]*>[^<]+")

	// 使用正则表达式匹配原始数据，获取到小说标题
	titleString := p.re.FindString(dataString)

	// 将 "<h1>《黑莲花女配重生了》" 根据 ">" 符号分割字符串，需要的是分割出来的第二项
	stringList := strings.Split(titleString, ">")

	novelTitle = stringList[1]

	// 去除可能存在的书名号
	novelTitle = strings.Trim(novelTitle, "《")
	novelTitle = strings.Trim(novelTitle, "》")

	// 记录日志
	zap.L().Info("[Parse Novel Title]: Successful parse of the novel title: " + novelTitle + ".")

	// 返回小说标题
	return novelTitle, nil
}

// parseNovelAuthor 解析小说作者
func (p *ParseQianyegeNovel) parseNovelAuthor(dataString string) (novelAuthor string, err error) {
	/*
		创建小说作者的正则表达式对象
			可以一次性获取到小说作者，且无需额外操作的正则表达式表达式：(?<=<p>作&nbsp;&nbsp;者：).*(?=</p>)
			Go 不支持这个写法！！！
	*/
	p.re = regexp.MustCompile("<p>作&nbsp;&nbsp;者：[^<]+")
	// 使用正则表达式匹配原始数据，获取到小说作者
	authorString := p.re.FindString(dataString)

	// "<p>作&nbsp;&nbsp;者：鹦鹉晒月" 根据 "：" 符号分割字符串，需要的是分割出来的第二项
	stringList := strings.Split(authorString, "：")

	novelAuthor = stringList[1]

	// 记录日志
	zap.L().Info("[Parse Novel Author]: Successful parse of the novel author: " + novelAuthor + ".")

	// 返回小说作者
	return novelAuthor, nil
}

// parseNovelAbstracts 解析小说简介
func (p *ParseQianyegeNovel) parseNovelAbstracts(dataString string) (novelAbstractsList []string, err error) {
	/*
		创建小说简介的正则表达式对象
			可以一次性获取到小说简介，且无需额外操作的正则表达式表达式：(?<=<div id="intro">\s*<p>&nbsp;&nbsp;&nbsp;&nbsp;).*(?=</p>\s*</div>)
			Go 不支持这个写法！！！
	*/
	p.re = regexp.MustCompile("&nbsp;&nbsp;&nbsp;&nbsp;[^/]+")

	// 使用正则表达式匹配原始数据，获取到小说简介
	novelAbstracts := p.re.FindString(dataString)

	novelAbstracts = strings.Trim(novelAbstracts, "<br")

	// 替换去除多余标签
	novelAbstracts = strings.ReplaceAll(novelAbstracts, "\r\n", "")
	novelAbstracts = strings.ReplaceAll(novelAbstracts, "<br>", "")
	novelAbstracts = strings.ReplaceAll(novelAbstracts, "<br/>", "")

	novelAbstractsList = strings.Split(novelAbstracts, "&nbsp;&nbsp;&nbsp;&nbsp;")

	// 记录日志
	zap.L().Info("[Parse Novel Abstracts]: Parse novel abstracts is success.")

	return novelAbstractsList, nil
}

// parseNovelGeneralChapters 解析小说的总章节
func (p *ParseQianyegeNovel) parseNovelGeneralChapters(getNovelMainUrl, dataString string) (generalChapter []*model.NovelGeneralChapter, err error) {
	// 创建网站域名的正则表达式对象
	p.re = regexp.MustCompile("^(http|https|ftp)://[a-zA-Z0-9-.]+")
	// 获取网站域名
	httpDomainName := p.re.FindString(getNovelMainUrl)

	/*
		创建小说总章节的正则表达式对象
			可以一次性获取到小说总章节，且无需额外操作的正则表达式表达式：(?<=<div id="intro">\s*<p>&nbsp;&nbsp;&nbsp;&nbsp;).*(?=</p>\s*</div>)
			Go 不支持这个写法！！！
	*/
	p.re = regexp.MustCompile("正文</dt>.*</dd>")

	// 正则表达式获取小说总章节部分
	generalChapterString := p.re.FindString(dataString)

	// 替换去除多余标签
	generalChapterString = strings.ReplaceAll(generalChapterString, "正文</dt>", "")
	generalChapterString = strings.ReplaceAll(generalChapterString, "<dd>", "")
	generalChapterString = strings.ReplaceAll(generalChapterString, "</dd>", "")

	/*
		现在的 generalChapterString 是如下模样：
			<a href="/64/64939/34154692.html">前引</a>
			<a href="/64/64939/34154693.html">001项家七小姐</a>
			<a href="/64/64939/34154694.html">002再相见</a>
			<a href="/64/64939/34154695.html">003忘初心</a>

		现在只需要按照 </a> 分割字符串，一个元素就是一个章节，注意：最后一个元素是空的
	*/
	generalChapterList := strings.Split(generalChapterString, "</a>")

	// 去除最后一个空元素
	generalChapterList = generalChapterList[:len(generalChapterList)-1]

	// 创建一个空的字典
	generalChapter = make([]*model.NovelGeneralChapter, 0)

	// 遍历每个章节
	for _, chapterString := range generalChapterList {
		// 获取章节标题位置
		chapterTitleIndex := strings.LastIndex(chapterString, ">")
		// 获取章节标题
		chapterTitle := chapterString[chapterTitleIndex+1:]

		// 获取章节链接Fast位置
		chapterLinkIndex := strings.Index(chapterString, "\"")
		// 获取章节链接Last位置
		chapterLinkLastIndex := strings.LastIndex(chapterString, "\"")
		// 获取章节链接
		chapterLink := chapterString[chapterLinkIndex+1 : chapterLinkLastIndex]

		// 将章节标题和章节链接存入字典中
		generalChapter = append(generalChapter, &model.NovelGeneralChapter{
			ChapterTitle: chapterTitle,
			ChapterURL:   httpDomainName + chapterLink,
		})
	}

	zap.L().Info("[Parse Novel General Chapters]: Successful parse of the entire novel, the novel a total of " + string(len(generalChapter)) + "chapters.")

	return generalChapter, nil
}
