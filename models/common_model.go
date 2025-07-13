package models

type CommentModel struct {
	Model
	Content      string        `gorm:"size:256" json:"content"`
	UserID       uint          `json:"userID"`
	UserModel    UserModel     `gorm:"ForeignKey:UserID" json:"-"`
	ArticleID    uint          `json:"articleID"`
	ArticleModel ArticleModel  `gorm:"ForeignKey:ArticleID" json:"-"`
	ParentID     *uint         `json:"parentID"` //父评论   允许为空，因为顶级评论没有父评论。
	ParentModel  *CommentModel `gorm:"ForeignKey:ParentID" json:"-"`
	//这个字段表示当前评论的所有子评论（回复）。
	//gorm:"ForeignKey:ParentID" 表示通过子评论的 ParentID 字段关联到当前评论的 ID（即当前评论是子评论的父评论）。
	SubCommentList []*CommentModel `gorm:"ForeignKey:ParentID" json:"-"`
	RootParentID   *uint           `json:"rootParentID"` //根评论
	DiggCount      int             `json:"diggCount"`    //评论点赞数

}
