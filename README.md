# Novel Downloader

简体中文 | [English](./README.en.md)

**novel-downloader** 是一个用 **Go** 语言编写的小说下载器，可以从能够访问的网站上将你喜欢的小说下载到本地，永久的将它们保存在你希望它存在的任何地方。

由于众所周知的原因，现在国内对小说有比较大的限制。在晋江、阅文之类的正版网站上购买的小说，可能因为种种原因，突然一天就无法观看了，而小说的版权却平台手里，平台不给看，作者因为没版权，也不能在其他地方补发。真金白银买的小说，就只能被永久的停留在记忆中，随着时间慢慢消失……若是有幸有些盗版网站恰好有你心心念念的小说都还好，若是小众一些的小说，那可能就回天乏术了。

基于以上原因，这个小说下载器就诞生了。

## 简单的使用方式

**novel-downloader** 目前属于初级阶段，交互方式还比较 **low**。

在`conf/`目录下的`config.toml`文件中：

1. 向`get_novel_main_url`项中添加需要下载的小说的主页链接
2. 向`save_novel_path`项中添加下载小说后，小说文件的保存路径（注意是**文件夹路径**！）

- 小说的主页链接：指的是包含小说的标题、作者、摘要、所有章节等信息的小说主页链接。

通常下载小说就只需要将`get_novel_main_url`项和`save_novel_path`项的信息填写好，直接运行程序就可以下载小说了。

## 扩展功能

**novel-downloader** 的主体功能是基于正则表达式实现的，由于不同的网站的网页结构不同，自然没办法单凭一套的代码来下载所有网站的小说。

所幸 **novel-downloader** 整体基于接口设计，可以很方便的扩展功能，可以很容易的扩展对其他网站的支持。

当需要对其他网站的小说支持时，只需要实现`internal/parsedata`目录下`parse_noveler.go`文件中定义的`ParseNoveler`接口。

若没有其他需要的话，其他的接口和方法都可以不改动，只要正确实现`ParseNoveler`接口，再在`internal/biz`的`biz_novel.go`文件中的`NewBizNovel()`函数中添加新定义的解析器，即可完成对新网站的支持。

`ParseNoveler`接口定义了 **novel-downloader** 中最最重要的功能模块，即解析小说的各个部分，包括：

- 公开方法：
   1. `ParseNovelMain()`：解析小说主页
      1. 使用私有方法`parseSourceDataEncode()`判断下载的源数据是否需要转换的编码方式，避免源数据乱码
      2. 将私有方法的2-5组合起来调用，完成小说主页的解析
   2. `ParseNovelChapters()`：解析小说章节
      1. 使用私有方法`parseSourceDataEncode()`判断下载的源数据是否需要转换的编码方式，避免源数据乱码
      2. 使用正则表达式解析小说章节的内容
- 私有方法：
   1. `parseSourceDataEncode()`：使用正则表达式解析源数据的编码
   2. `parseNovelTitle()`：使用正则表达式解析小说标题
   3. `parseNovelAuthor()`：使用正则表达式解析小说作者
   4. `parseNovelAbstracts()`：使用正则表达式解析小说简介
   5. `parseNovelGeneralChapters()`：使用正则表达式解析小说总章节

