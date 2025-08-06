package desensitization

import "strings"

func TelDesensitization(tel string) string {
	if len(tel) != 11 {
		return ""
	}
	// 17564263096
	// 175****3096
	return tel[:3] + "****" + tel[7:]
}

func EmailDesensitization(email string) string {
	if email == "" {
		return ""
	}
	slice := strings.Split(email, "@")
	//154188888@qq.com
	//1****@qq.com
	return slice[0][:1] + "****@" + slice[1]

}
