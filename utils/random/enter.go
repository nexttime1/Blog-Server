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
