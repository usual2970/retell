package encode

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// ToPinyin 将汉字转换为拼音
// 添加一个判断字符是否为汉字的辅助函数
func isHanzi(r rune) bool {
	return r >= '\u4e00' && r <= '\u9fff'
}

// 添加一个检查字符串是否包含汉字的辅助函数
func containsHanzi(text string) bool {
	for _, r := range text {
		if isHanzi(r) {
			return true
		}
	}
	return false
}

func ToPinyin(text string) string {
	// 检查是否包含汉字
	if !containsHanzi(text) {
		rs := strings.ToLower(text)
		return strings.ReplaceAll(rs, "_", "")
	}

	args := pinyin.NewArgs()
	args.Style = pinyin.Normal

	// 转换为拼音
	pinyinSlice := pinyin.Pinyin(text, args)

	// 将二维切片转换为一维
	result := make([]string, 0, len(pinyinSlice))
	for _, p := range pinyinSlice {
		if len(p) > 0 {
			result = append(result, p[0])
		}
	}

	rs := strings.Join(result, "")
	return strings.ReplaceAll(rs, "_", "")
}
