package main

import (
	"encoding/json"
	"fmt"
	"github.com/liu-cn/json-filter/filter"
	"time"
)

type User struct {
	UID        uint      `json:"uid,select(article)"`
	Avatar     string    `json:"avatar,select(article)"`
	Nickname   string    `json:"nickname,select(article|profile|list)"`
	Sex        int       `json:"sex,select(profile)"`
	VipEndTime time.Time `json:"vip_end_time,select(profile),omit(vip)"`
	Price      string    `json:"price,select(profile)"`
}

func main() {

	user := User{
		UID:        1,
		Nickname:   "boyan",
		Avatar:     "avatar",
		Sex:        1,
		VipEndTime: time.Now().Add(time.Hour * 24 * 365),
		Price:      "999.9",
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal)) //The following is the official JSON parsing output: you can see that all fields have been parsed
	//{"uid":1,"nickname":"boyan","avatar":"avatar","sex":1,"vip_end_time":"2023-03-06T23:11:22.622693+08:00","price":"999.9"}

	//usage：filter.Select("select case",This can be：slice/array/struct/pointer/map)
	fmt.Println(filter.Select("article", user)) //The following is the JSON filtered by JSON filter. This output is the JSON under the article interface
	//{"avatar":"avatar","nickname":"boyan","uid":1}

	fmt.Println(filter.Select("profile", user)) //profile result
	//{"nickname":"boyan","price":"999.9","sex":1,"vip_end_time":"2023-03-06T23:31:28.636529+08:00"}
	fmt.Println(filter.Select("list", user))
	fmt.Println(filter.Omit("vip", user)) //过滤完 vip
}
