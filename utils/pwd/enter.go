package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPwd 将明文密码转为 hash密码
func HashPwd(pwd string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// CheckPwd 对比hash密码  与明文密码是否为一对
func CheckPwd(hashPwd string, pwd string) bool {
	byteHash := []byte(hashPwd)
	bytePwd := []byte(pwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
