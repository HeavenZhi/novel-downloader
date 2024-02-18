package encode

type Encoder interface {
	// 转换为 GBK 编码
	EncodeToGBK(data []byte) ([]byte, error)
	// 转换为 UTF-8 编码
	EncodeToUTF8(data []byte) ([]byte, error)
}
