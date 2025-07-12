// testdata/10.获取地理位置.go
package main

import (
	"Blog_server/core"
	"fmt"
)

func main() {
	core.InitIPDB()
	addr := core.GetIpAddr("112.224.194.255")
	fmt.Println(addr)
	fmt.Println(core.GetIpAddr("162.224.194.255"))
	fmt.Println(core.GetIpAddr("10.224.194.255"))
	fmt.Println(core.GetIpAddr("1022.224.194.255"))
}
