package common

import (
	"Blog_server/conf"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// ToYAML 将配置结构体写入 YAML 文件
// 参数：
//   - filePath: 要写入的 YAML 文件路径
//   - config: 要写入的配置结构体实例
//
// 返回：
//   - 错误信息，如果成功则返回 nil
func ToYAML(filePath string, config *conf.Config) error {
	// 1. 将配置结构体序列化为 YAML 格式字节流
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("YAML 序列化失败: %w", err)
	}

	// 2. 打开文件（不存在则创建，存在则覆盖）
	// 权限 0644 表示：所有者可读写，组和其他用户可读
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close() // 确保文件最终关闭

	// 3. 将 YAML 字节流写入文件
	_, err = file.Write(yamlData)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
