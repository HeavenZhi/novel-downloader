package getdata

import "net/http"

type GetNoveler interface {
	// 获取小说数据
	GetNovelData(downloadURL string, proxyList, userAgentList []string) (data []byte, err error)
	// 根据代理服务器情况，选择性创建 http 客户端
	createHttpClient(proxyList []string) (httpClient *http.Client, err error)
	// 根据 User-Agent 情况，选择性创建 http 请求
	createHttpRequest(downloadURL string, userAgentList []string) (request *http.Request, err error)
}
