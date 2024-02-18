package config

// Configer 配置接口，方便切换配置文件框架
type Configer interface {
	Init(filePath string) (err error)
}

// Conf 定义全局变量
var Conf = new(Config)

/*
使用 viper.GetXxx() 读取方式
注意：
Viper 使用的是 mapstructure 为标签
*/

// Config 配置结构体
type Config struct {
	*NovelConfig             `mapstructure:"novel"`
	*SettingConfig           `mapstructure:"setting"`
	*SaveMarkdownStyleConfig `mapstructure:"save_markdown_style"`
	*SaveTxtStyleConfig      `mapstructure:"save_txt_style"`
	*LogConfig               `mapstructure:"log"`
}

// NovelConfig 小说配置结构体
type NovelConfig struct {
	//小说主页地址
	GetNovelMainUrl string `mapstructure:"get_novel_main_url"`
	// 小说保存位置
	SaveNovelPath string `mapstructure:"save_novel_path"`
}

// NovelSettingConfig 小说下载器结构体
type SettingConfig struct {
	// 下载小说时单个 goroutine 的休眠时间（防止被封IP,最低设置为3秒）
	DownloadSleepTime int64 `mapstructure:"download_sleep_time"`
	// Proxy 代理服务器池，不断切换随机IP访问（避免被目标网站封IP）
	ProxyList []string `mapstructure:"proxy_list"`
	// User-Agent 池，不断切换随机User-Agent访问（避免被目标网站封IP）
	UserAgentList []string `mapstructure:"user_agent_list"`
}

// SaveMarkdownStyleConfig 小说 markdown 保存风格结构体
type SaveMarkdownStyleConfig struct {
	// 首行缩进风格
	IndentStyle string `mapstructure:"indent_style"`
	// 换行符风格
	NewLineStyle string `mapstructure:"new_line_style"`
	// 小说源展示风格
	NovelSoureStyle string `mapstructure:"novel_soure_style"`
	// 小说标题左边展示风格
	NovelTitleLeftStyle string `mapstructure:"novel_title_left_style"`
	// 小说标题右边展示风格
	NovelTitleRightStyle string `mapstructure:"novel_title_right_style"`
	// 小说作者展示风格
	NovelAuthorStyle string `mapstructure:"novel_author_style"`
	// 小说简介展示风格
	NovelAbstractsStyle string `mapstructure:"novel_abstracts_style"`
	// 标题分隔符风格
	TitleSeparatorStyle string `mapstructure:"title_separator_style"`
}

// SaveTxtStyleConfig 小说 txt 保存风格结构体
type SaveTxtStyleConfig struct {
	// 首行缩进风格
	IndentStyle string `mapstructure:"indent_style"`
	// 换行符风格
	NewLineStyle string `mapstructure:"new_line_style"`
	// 小说源展示风格
	NovelSoureStyle string `mapstructure:"novel_soure_style"`
	// 小说标题左边展示风格
	NovelTitleLeftStyle string `mapstructure:"novel_title_left_style"`
	// 小说标题右边展示风格
	NovelTitleRightStyle string `mapstructure:"novel_title_right_style"`
	// 小说作者展示风格
	NovelAuthorStyle string `mapstructure:"novel_author_style"`
	// 小说简介展示风格
	NovelAbstractsStyle string `mapstructure:"novel_abstracts_style"`
	// 标题分隔符风格
	TitleSeparatorStyle string `mapstructure:"title_separator_style"`
}

// LogConfig 日志配置结构体
type LogConfig struct {
	Level         string `mapstructure:"level"`
	Filename      string `mapstructure:"filename"`
	ErrorFilename string `mapstructure:"error_filename"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxAge        int    `mapstructure:"max_age"`
	MaxBackups    int    `mapstructure:"max_backups"`
}
