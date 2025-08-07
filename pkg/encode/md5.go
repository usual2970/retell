package encode

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 将给定字符串转换为32位小写MD5字符串
func Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
