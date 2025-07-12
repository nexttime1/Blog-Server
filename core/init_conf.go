package core

import (
	"Blog_server/conf"
	"Blog_server/flags"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// 从main文件的 根目录

func ReadConf() *conf.Config {
	file, err := os.ReadFile(flags.FileOption.File)
	if err != nil {
		panic(err)
	}
	var c = new(conf.Config)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件格式错误 ,%s", err))
	}

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)
	fmt.Println(c)

	return c
}
