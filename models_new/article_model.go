package models_new

type ArticleModel struct {
	Model
	Title      string `gorm:"size:32" json:"title"`
	Abstract   string `gorm:"size:256" json:"abstract"`
	Content    string `json:"content"`
	CategoryID string `json:"categoryID"`

	CommentModels []CommentModel `gorm:"foreignKey:ArticleID" json:"-"`                // 文章的评论列表
	TagList       string         `gorm:"type:longtext;serializer:json" json:"tagList"` //标签列表
	Cover         string         `gorm:"size:256" json:"cover"`                        //封面
	UserID        uint           `json:"userID"`
	UserModel     UserModel      `gorm:"ForeignKey:UserID" json:"-"`
	LookCount     int            `json:"lookCount"`
	DiggCount     int            `json:"diggCount"`    //点赞
	CommentCount  int            `json:"commentCount"` //评论
	CollectsCount int            `json:"collectCount"` // 收藏
	OpenCount     bool           `json:"openCount"`    //是否开启评论         标记无
	Status        int8           `json:"status"`       //状态   草稿  审核中  已发布   标记无

	Category   string      `gorm:"size:20" json:"category"`      // 文章分类
	Source     string      `json:"source"`                       // 文章来源
	Link       string      `json:"link"`                         // 原文链接
	Banner     BannerModel `gorm:"foreignKey:BannerID" json:"-"` // 文章封面
	BannerID   uint        `json:"banner_id"`                    // 文章封面id
	NickName   string      `gorm:"size:42" json:"nick_name"`     // 发布文章的用户昵称
	BannerPath string      `json:"banner_path"`                  // 文章的封面
	//Tags       .Array `gorm:"type:string;size:64" json:"tags"`
}
