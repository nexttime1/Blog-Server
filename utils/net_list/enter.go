package net_list

import (
	"github.com/sirupsen/logrus"
	"net"
)

func GetIpList() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		logrus.Errorf("%s", err)
		return nil
	}
	var ipList []string

	for _, i := range interfaces {
		address, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range address {
			ipNet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			ipList = append(ipList, ip.String())
		}

	}
	return ipList
}
