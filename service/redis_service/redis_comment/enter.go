package redis_comment

import (
	"Blog_server/global"
	"fmt"
	"strconv"
)

const CommentPrefix = "comment"

func Comment(id string) {
	num, err := global.Redis.HGet(CommentPrefix, id).Int()
	fmt.Errorf("%s", err)
	num++
	global.Redis.HSet(CommentPrefix, id, num)
}

func GetComment(id string) int {
	num, _ := global.Redis.HGet(CommentPrefix, id).Int()
	return num
}

func GetCommentInfo() map[string]int {
	var CommentInfo = make(map[string]int)
	data := global.Redis.HGetAll(CommentPrefix).Val()
	for id, val := range data {
		num, _ := strconv.Atoi(val)
		CommentInfo[id] = num
	}
	return CommentInfo
}

func CommentInfoClear() {
	global.Redis.Del(CommentPrefix)
}
