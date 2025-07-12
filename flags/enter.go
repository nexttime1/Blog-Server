package flags

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FileOption = new(Options)

func Parse() {
	flag.StringVar(&FileOption.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FileOption.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FileOption.Version, "v", false, "版本")
	flag.Parse()
}

func Run() {
	if FileOption.DB { //数据库迁移
		fmt.Println("数据库开始迁移")
		FlagDB()
		os.Exit(0)
	}

}
