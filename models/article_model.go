package models

type ArticleModel struct {
	Model
	Title        string    `gorm:"size:32" json:"title"`
	Abstract     string    `gorm:"size:256" json:"abstract"`
	Content      string    `json:"content"`
	CategoryID   string    `json:"categoryID"`
	TagList      string    `gorm:"type:longtext;serializer:json" json:"tagList"` //标签列表
	Cover        string    `gorm:"size:256" json:"cover"`                        //封面
	UserID       uint      `json:"userID"`
	UserModel    UserModel `gorm:"ForeignKey:UserID" json:"-"`
	LookCount    int       `json:"lookCount"`
	DiggCount    int       `json:"diggCount"`    //点赞
	CommentCount int       `json:"commentCount"` //评论
	CollectCount int       `json:"collectCount"` // 收藏
	OpenCount    bool      `json:"openCount"`    //是否开启评论
	Status       int8      `json:"status"`       //状态   草稿  审核中  已发布

}
