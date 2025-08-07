package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func StringToMD5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func RemoveTagetStr(str, sub string) string {
	if str == "" || sub == "" {
		return ""
	}
	if strings.Contains(str, sub) {
		// 删除子字符串
		str = strings.Replace(str, sub, "", -1)
	}
	return str
}
