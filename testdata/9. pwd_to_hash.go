package main

import (
	"Blog_server/utils/pwd"
	"fmt"
)

func main() {
	password := "123456"
	//hashPwd := pwd.HashPwd(password)

	ok := pwd.CheckPwd("$2a$04$Yne3RA63lj82zZI/O35X5OyrBFG19UTBO6eTL.XWLBWK4WnDEy2HK", password)
	fmt.Println(ok)

}
