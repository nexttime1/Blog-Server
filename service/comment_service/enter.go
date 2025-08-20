package comment_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/service/redis_service/redis_count"
	"Blog_server/utils/jwts"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CommentAddRequest struct {
	ArticleID       string `json:"article_id" binding:"required" msg:"请选择文章"`
	Content         string `json:"content" binding:"required" msg:"请输入评论内容"`
	ParentCommentID *uint  `json:"parent_comment_id"` // 父评论id
}

type CommentDeleteRequest struct {
	ID int `form:"id" uri:"id"`
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

func CommentDeleteService(cr CommentDeleteRequest) (int, error) {

	//判断 评论存不存在
	var commentModel models.CommentModel
	err := global.DB.Where("id = ?", cr.ID).Take(&commentModel).Error
	if err != nil {
		logrus.Errorf("评论不存在  %s", err)
		return 0, fmt.Errorf("评论不存在 %s", err)
	}
	// 文章一定存在  看看是不是子评论
	if commentModel.ParentCommentID != nil {
		//子评论 让父评论评论数 -1
		var ParentsModel models.CommentModel
		global.DB.Where("id = ?", commentModel.ParentCommentID).Take(&ParentsModel)
		global.DB.Model(&ParentsModel).Update("comment_count", gorm.Expr("comment_count - ?", 1))
	}
	// 找出他们子评论的id
	var subCommentModels []*models.CommentModel
	inList := FindSubModel(commentModel, &subCommentModels)
	//加上自己删掉
	inList = append(inList, uint(cr.ID))
	global.DB.Model(models.CommentModel{}).Delete("id in ?", inList)
	//删除redis上的 文章评论数
	count := len(inList)
	redis_count.NewComment().SetNum(commentModel.ArticleID, -count)
	// 所有评论点赞删除掉
	for _, num := range inList {
		id := fmt.Sprintf("%d", num)
		redis_count.NewCommentDigg().SetNum(id, 0)
	}

	return count, err

}
func FindSubModel(model models.CommentModel, subCommentModels *[]*models.CommentModel) []uint {
	var idList []uint
	global.DB.Preload("SubComments").Take(&model)

	for _, commentModel := range model.SubComments {
		idList = append(idList, commentModel.ID)
		*subCommentModels = append(*subCommentModels, commentModel)
		FindSubModel(*commentModel, subCommentModels)
		commentModel.SubComments = nil
	}
	return idList
}
