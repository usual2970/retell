package encode

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rng  *rand.Rand
	once sync.Once
)

func init() {
	once.Do(func() {
		// 使用当前时间作为种子创建随机数生成器
		source := rand.NewSource(time.Now().UnixNano())
		rng = rand.New(source)
	})
}

func Uri(prefix ...string) string {
	pre := ""
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	return pre + time.Now().Format("0601021504") + GenerateRandomString(6)
}

func AppId() string {
	rs := Md5(GenerateRandomString(32))
	return "yu3" + rs[:18]
}

func GenerateRandomString(length int) string {
	// 定义字符集
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成随机字符串
	randomString := make([]byte, length)
	for i := range length {
		randomString[i] = charset[rng.Intn(len(charset))]
	}

	return string(randomString)
}
