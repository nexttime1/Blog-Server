package random

import (
	"fmt"
	"math/rand"
	"time"
)

// DigitCode 生成4位随机验证码
func DigitCode() string {
	// 创建新的随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成1000到9999之间的随机数（包含1000和9999）
	number := r.Intn(9000) + 1000

	return fmt.Sprintf("%d", number)
}

// GenerateRandomString 生成一个随机字符串
func GenerateRandomString(count int) string {
	// 初始化随机数
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 字符集：小写字母+数字
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	// 结果字节切片
	b := make([]byte, count)

	// 生成每个字符
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
