package core

import (
	"Blog_server/global"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func InitEs() *elastic.Client {
	//创建客户端实例   把嗅探功能关闭
	es := global.Config.Es
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(es.Username, es.Password))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	logrus.Infof("ES 连接成功")
	return client
}
