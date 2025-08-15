package redis_look

import (
	"Blog_server/global"
	"fmt"
	"strconv"
)

const LookPrefix = "look"

func Look(id string) {
	num, err := global.Redis.HGet(LookPrefix, id).Int()
	fmt.Errorf("%s", err)
	num++
	global.Redis.HSet(LookPrefix, id, num)

}

func GetLook(id string) int {
	num, _ := global.Redis.HGet(LookPrefix, id).Int()
	return num
}

func GetLookInfo() map[string]int {
	var LookInfo = make(map[string]int)
	data := global.Redis.HGetAll(LookPrefix).Val()
	for id, val := range data {
		num, _ := strconv.Atoi(val)
		LookInfo[id] = num
	}
	return LookInfo
}

func LookClear() {
	global.Redis.Del(LookPrefix)
}
