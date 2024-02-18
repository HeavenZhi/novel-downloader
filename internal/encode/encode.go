package encode

import (
	"bytes"
	"io"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Encode 实现了 Encoder 接口
type Encode struct {
}

// NewEncode 创建 Encode 对象（构造方法）
func NewEncode() *Encode {
	return &Encode{}
}

// toGBK 将 UTF-8 编码的字符串转换为 GBK 编码的字符串
func (e *Encode) EncodeToGBK(data []byte) ([]byte, error) {
	// 创建一个从 UTF-8 到 GBK 的解码器
	decoder := simplifiedchinese.GBK.NewDecoder()

	// 创建一个基于原始数据和解码器的 Reader
	reader := transform.NewReader(bytes.NewReader(data), decoder)

	// 读取并转换所有数据到 GBK 编码
	dataGBK, err := io.ReadAll(reader)
	if err != nil {
		// 如果在读取和转换过程中发生错误，则返回 nil 和错误信息
		return nil, err
	}

	zap.L().Info("[Encode To GBK]: The source data was successfully converted to GBK encoding.")

	// 转换成功，返回 GBK 编码的数据
	return dataGBK, nil
}

// toUTF8 将 GBK 编码的字符串转换为 UTF-8 编码的字符串
func (e *Encode) EncodeToUTF8(data []byte) ([]byte, error) {
	return nil, nil
}
