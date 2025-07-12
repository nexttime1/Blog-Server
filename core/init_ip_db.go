package core

import (
	UtilsIp "Blog_server/utils/ip"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
	"strings"
)

var searcher *xdb.Searcher

func InitIPDB() {
	var dbPath = "init/ip2region.xdb"
	_searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		logrus.Errorf("初始化失败 %s", err.Error())
		return
	}
	searcher = _searcher

}

func GetIpAddr(ip string) (ipAddr string) {
	addr := UtilsIp.HasLocalIPAddr(ip)
	if addr {
		return "内网ip"
	}

	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logrus.Warnf(" ip地址错误  %s", err.Error())
		return "错误的ip地址"
	}

	ListDdata := strings.Split(region, "|") //中国|0|山东省|济南市|联通
	if len(ListDdata) != 5 {
		logrus.Warnf("异常的Ip地址", ip)
		return "未知地址"

	}
	// 国家 0 省份 市区 运营商
	country := ListDdata[0]
	province := ListDdata[2]
	city := ListDdata[3]

	if province != "0" && city != "0" {
		return fmt.Sprintf("%s %s", province, city)
	}
	if country != "0" && province != "0" {
		return fmt.Sprintf("%s %s", country, province)
	}
	if country != "0" { // 新加坡|0|0|0|xx
		return fmt.Sprintf("%s", country)
	}

	return region
}
