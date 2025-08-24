package net_ip

import (
	"Blog_server/global"
	"Blog_server/utils/net_list"
	"github.com/sirupsen/logrus"
)

func PrintSystem() {
	ip := global.Config.System.IP
	port := global.Config.System.Port

	if ip == "0.0.0.0" {
		iplist := net_list.GetIpList()
		for _, newIp := range iplist {
			logrus.Infof("gvb_server 运行在：http://%s:%d/api", newIp, port)
			logrus.Infof("gvb_server api 文档 运行在：http://%s:%d/swagger/index.html#", newIp, port)

		}
	} else {
		logrus.Infof("gvb_server 运行在：http://%s:%d/api", ip, port)
		logrus.Infof("gvb_server api 文档 运行在：http://%s:%d/swagger/index.html#", ip, port)
	}

}
