package getdata

import (
	"crypto/tls"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// GetNovel 实现 GetNoveler 接口
type GetNovel struct {
}

func NewGetNovel() *GetNovel {
	return &GetNovel{}
}

// GetNovelData 获取小说数据
func (g *GetNovel) GetNovelData(downloadURL string, proxyList, userAgentList []string) (data []byte, err error) {
	// 根据代理服务器情况，选择性创建 http 客户端
	httpClient, err := g.createHttpClient(proxyList)
	if err != nil {
		return nil, err
	}

	// 根据 User-Agent 情况，选择性创建 http 请求
	request, err := g.createHttpRequest(downloadURL, userAgentList)
	if err != nil {
		return nil, err
	}

	// 发送请求
	response, err := httpClient.Do(request)
	if err != nil {
		zap.L().Error("[Error Get Novel Data]: Send request error:", zap.Error(err))
		return nil, err
	}

	// defer 延迟执行，函数结束后关闭相关链接
	defer response.Body.Close()

	// 读取响应内容
	data, err = io.ReadAll(response.Body)
	if err != nil {
		zap.L().Error("[Error Get Novel Data]: Read response error:", zap.Error(err))
		return nil, err
	}

	// 判断小说主页数据是否为空
	if len(data) == 0 {
		zap.L().Error("[Error Get Novel Data]: No novel data was obtained. Please check the novel home page link: " + downloadURL)
		return nil, err
	}

	// 向日志中记录数据
	zap.L().Info("[Get Novel Data]: Get novel data success.")

	return data, nil
}

// createHttpClient 根据代理服务器情况，选择性创建 http 客户端
func (g *GetNovel) createHttpClient(proxyList []string) (httpClient *http.Client, err error) {
	// 获取代理服务器数量
	randProxy := len(proxyList)

	// 当存在代理服务器时，使用随机代理服务器创建 http 客户端，连接目标网站
	if randProxy > 0 {
		// 初始化随机种子数
		rand.Seed(time.Now().Unix())

		// 选中随机数
		proxy := rand.Intn(randProxy)

		// 根据随机的代理服务器 string 创建 代理服务器的 url 对象
		proxyURL, err := url.Parse(proxyList[proxy])
		if err != nil {
			zap.L().Error("[Error Get Novel Data]: Parse url error:", zap.Error(err))
			return nil, err
		}

		// 创建 http 客户端
		httpClient = &http.Client{
			Transport: &http.Transport{
				// 使用随机的代理服务器访问目标网站
				Proxy: http.ProxyURL(proxyURL),
				// 跳过证书验证
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		zap.L().Info("[Create HTTP Client]: Create http client success(Use proxy: " + proxyList[proxy] + ").")

	} else {
		// 不存在代理服务器时，直接创建 http 客户端
		httpClient = &http.Client{
			Transport: &http.Transport{
				// 跳过证书验证
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		zap.L().Info("[Create HTTP Client]: Create http client success(No proxy use).")
	}

	return httpClient, nil
}

// createHttpRequest 根据 User-Agent 情况，选择性创建 http 请求
func (g *GetNovel) createHttpRequest(downloadURL string, userAgentList []string) (request *http.Request, err error) {
	// 创建请求
	request, err = http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		zap.L().Error("[Error Get Novel Data]: Create request error: ", zap.Error(err))
		return nil, err
	}

	// 自定义 Header，保持连接
	request.Header.Set("Connection", "keep-alive")

	// 获取 User-Agent 数量
	randUserAgent := len(userAgentList)

	// 当存在 User-Agent 时，使用随机 User-Agent 发送请求
	if randUserAgent > 0 {
		// 初始化随机种子数
		rand.Seed(time.Now().Unix())

		// 选中随机数
		userAgent := rand.Intn(randUserAgent)

		// 自定义 Header，使用随机的 User-Agent 发送请求
		request.Header.Set("User-Agent", userAgentList[userAgent])

		zap.L().Info("[Create HTTP Request]: Create http request success(Use User-Agent: " + userAgentList[userAgent] + ").")
	} else {
		// 不存在 User-Agent 时，自定义 Header 固定的 User-Agent 发送请求
		request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")

		zap.L().Info("[Create HTTP Request]: Create http request success(No User-Agent use).")
	}

	return request, nil
}
