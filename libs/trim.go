package libs

import "strings"

// 对剪贴板的字符串进行处理
func DefaultTrimmer(payload string) (result string) {
	// 去除最后一个回车.
	result = strings.TrimSuffix(payload, "\n") // 和 TrimRight 区别
	return
}
