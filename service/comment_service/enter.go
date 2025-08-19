package comment_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/redis_service/redis_count"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/sirupsen/logrus"
)

type CommentAddRequest struct {
	ArticleID       string `json:"article_id" binding:"required" msg:"请选择文章"`
	Content         string `json:"content" binding:"required" msg:"请输入评论内容"`
	ParentCommentID *uint  `json:"parent_comment_id"` // 父评论id
}

func CommentAddService(cr CommentAddRequest, claim *jwts.MyClaims) error {
	//判断用户存不存在
	var user models.UserModel
	err := global.DB.Where("id = ?", claim.UserID).Take(&user).Error
	if err != nil {
		logrus.Errorf("用户不存在  %s", err)
		return fmt.Errorf("用户不存在")
	}
	// 判断文章存不存在
	var article models.ArticleModel
	err = article.ExistById(cr.ArticleID)
	if err != nil {
		logrus.Errorf("文章不存在  %s", err)
		return fmt.Errorf("文章不存在")
	}

	// 判断父评论
	if cr.ParentCommentID != nil {
		//我的父评论  一定 和我一个文章  不是一个文章不行
		var commentModel models.CommentModel
		err = global.DB.Where("id = ?", cr.ParentCommentID).Take(&commentModel).Error
		if err != nil {
			logrus.Errorf("父评论不存在  %s", err)
			return fmt.Errorf("父评论不存在")
		}
		if commentModel.ArticleID != cr.ArticleID {
			logrus.Errorf("父评论和自己不是一篇文章")
			return fmt.Errorf("父评论和自己不是一篇文章")
		}
		//给父评论数 + 1
		global.DB.Model(&commentModel).Update("comment_count", commentModel.CommentCount+1)
	}
	// 创建评论
	err = global.DB.Create(&models.CommentModel{
		ParentCommentID: cr.ParentCommentID,
		Content:         cr.Content,
		ArticleID:       cr.ArticleID,
		UserID:          claim.UserID,
	}).Error
	if err != nil {
		logrus.Errorf("创建评论失败  %s", err)
		return fmt.Errorf("创建评论失败")
	}
	//文章评论 +1
	redis_count.NewComment().Set(cr.ArticleID)
	return nil
}
