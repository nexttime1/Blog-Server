package ip

import "net"

//三判断   第一 环回地址  第二 是不是ipv4  第三 是不是 内网指定

func HasLocalIPAddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip)) // ParseIP  将字符串解析为 net.IP 类型
}

// HasLocalIP 检测 IP 地址是否是内网地址
// 通过直接对比ip段范围效率更高
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() { //  回环地址检查  127
		return true
	}

	ip4 := ip.To4() //  转成ipv4  若失败（如 IPv6 地址），返回 nil
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}
