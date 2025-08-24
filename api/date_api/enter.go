package date_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"time"
)

type DateApi struct {
}

type DateCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DateResponse struct {
	DateList  []string
	LoginData []int
	SignData  []int
}

type DataSumResponse struct {
	UserCount      int `json:"user_count"`
	ArticleCount   int `json:"article_count"`
	MessageCount   int `json:"message_count"`
	ChatGroupCount int `json:"chat_group_count"`
	NowLoginCount  int `json:"now_login_count"`
	NowSignCount   int `json:"now_sign_count"`
}

// SevenLoginView 七天之内登录情况
// @Summary 七天内情况
// @Description 七天之内登录情况，注册情况
// @Tags 数据统计
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 400 {object} res.Response "请求参数错误"
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/data_login [get]
func (DateApi) SevenLoginView(c *gin.Context) {
	var loginDateCount, signDateCount []DateCount

	global.DB.Model(models.LoginDataModel{}).
		Where("date_sub(curdate(), interval 7 day) <= created_at").
		Select("date_format(created_at, '%Y-%m-%d') as date", "count(id) as count").
		Group("date").
		Scan(&loginDateCount)
	global.DB.Model(models.UserModel{}).
		Where("date_sub(curdate(), interval 7 day) <= created_at").
		Select("date_format(created_at, '%Y-%m-%d') as date", "count(id) as count").
		Group("date").
		Scan(&signDateCount)

	var LoginDateMap = make(map[string]int)
	var SignDateMap = make(map[string]int)
	for _, v := range loginDateCount {
		LoginDateMap[v.Date] = v.Count
	}
	for _, v := range signDateCount {
		SignDateMap[v.Date] = v.Count
	}
	var DateList = make([]string, 0)
	var LoginDateList = make([]int, 0)
	var SignDateList = make([]int, 0)
	now := time.Now()
	for i := -6; i <= 0; i++ {
		day := now.AddDate(0, 0, i).Format("2006-01-02")
		DateList = append(DateList, day)
		LoginDateList = append(LoginDateList, LoginDateMap[day])
		SignDateList = append(SignDateList, SignDateMap[day])
	}

	res.OkWithData(c, DateResponse{
		DateList:  DateList,
		LoginData: LoginDateList,
		SignData:  SignDateList,
	})
}

// DataSumView 数据统计
// @Summary 统计各字段
// @Description 统计各字段
// @Tags 数据统计
// @Param token header string true "用户认证令牌"
// @Success 200 {object} res.Response{}
// @Failure 500 {object} res.Response "服务器内部错误"
// @Router /api/data_sum [get]
func (DateApi) DataSumView(c *gin.Context) {
	var userCount, articleCount, messageCount, ChatGroupCount int
	var nowLoginCount, nowSignCount int

	//文章总数
	result, _ := global.Es.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewMatchAllQuery()).
		Do(context.Background())

	articleCount = int(result.Hits.TotalHits.Value)
	//总用户
	global.DB.Model(models.UserModel{}).Select("count(id)").Scan(&userCount)
	//总消息
	global.DB.Model(models.MessageModel{}).Select("count(id)").Scan(&messageCount)

	//聊天室群发消息
	global.DB.Model(models.ChatModel{IsGroup: true}).Select("count(id)").Scan(&ChatGroupCount)

	//今天登录人数
	global.DB.Model(models.LoginDataModel{}).Where("to_days(created_at)=to_days(now())").
		Select("count(id)").Scan(&nowLoginCount)
	//今天注册人数
	global.DB.Model(models.UserModel{}).Where("to_days(created_at)=to_days(now())").
		Select("count(id)").Scan(&nowSignCount)

	res.OkWithData(c, DataSumResponse{
		UserCount:      userCount,
		ArticleCount:   articleCount,
		MessageCount:   messageCount,
		ChatGroupCount: ChatGroupCount,
		NowLoginCount:  nowLoginCount,
		NowSignCount:   nowSignCount,
	})

}
