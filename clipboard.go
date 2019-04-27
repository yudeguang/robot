package robot

import (
	"github.com/atotto/clipboard"
)

//设置文本到剪切板,设置为nil则表示清空剪切板
func SetClipboard(text string) error {
	return clipboard.WriteAll(text)
}

//从剪切板中获取文本数据
func GetClipboard() (string, error) {
	return clipboard.ReadAll()
}
