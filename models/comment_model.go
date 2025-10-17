package models

import (
	"Blog_server/global"
	"gorm.io/gorm"
)

// CommentModel 评论表
type CommentModel struct {
	Model              `json:",select(c)"`
	SubComments        []*CommentModel `gorm:"foreignkey:ParentCommentID" json:"sub_comments,select(c)"` // 子评论列表
	ParentCommentModel *CommentModel   `gorm:"foreignkey:ParentCommentID" json:"comment_model"`          // 父级评论
	ParentCommentID    *uint           `json:"parent_comment_id,select(c)"`                              // 父评论id
	Content            string          `gorm:"size:256" json:"content,select(c)"`                        // 评论内容
	DiggCount          int             `gorm:"size:8;default:0;" json:"digg_count,select(c)"`            // 点赞数
	CommentCount       int             `gorm:"size:8;default:0;" json:"comment_count,select(c)"`         // 子评论数
	ArticleID          string          `gorm:"size:32" json:"article_id,select(c)"`                      // 文章id
	User               UserModel       `json:"user,select(c)"`                                           // 关联的用户
	UserID             uint            `json:"user_id,select(c)"`                                        // 评论的用户
}

func (c *CommentModel) BeforeDelete(tx *gorm.DB) (err error) {
	// 先把子评论删掉
	return nil
}

// FindAllSubCommentList 找一个评论的所有子评论,一维化
func FindAllSubCommentList(com CommentModel) (subList []CommentModel) {
	global.DB.Preload("SubComments").Preload("User").Take(&com)
	for _, model := range com.SubComments {
		subList = append(subList, *model)
		subList = append(subList, FindAllSubCommentList(*model)...)
	}
	return
}

// GetCommentTree 获取评论树
func GetCommentTree(rootComment *CommentModel) *CommentModel {
	global.DB.Preload("User").Preload("SubComments").Find(rootComment)

	// 递归获取子评论树
	for _, subComment := range rootComment.SubComments {
		GetCommentTree(subComment)
	}

	return rootComment
}
