package core

import (
	"Blog_server/flags"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// 从main文件的 根目录
var confPath = "settings.yaml"

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}
type Config struct {
	System System `yaml:"system"`
}

func ReadConf() {
	file, err := os.ReadFile(flags.FileOption.File)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件格式错误 ,%s", err))
	}

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)
}
