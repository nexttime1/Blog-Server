package article_service

import (
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/utils/jwts"
	"errors"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"math/rand"
	"strings"
	"time"
)

type ArticleAddRequest struct {
	Title    string     `json:"title" binding:"required" msg:"文章标题必填"`   // 文章标题
	Abstract string     `json:"abstract"`                                      // 文章简介
	Content  string     `json:"content" binding:"required" msg:"文章内容必填"` // 文章内容
	Category string     `json:"category"`                                      // 文章分类
	Source   string     `json:"source"`                                        // 文章来源
	Link     string     `json:"link"`                                          // 原文链接
	BannerID uint       `json:"banner_id"`                                     // 文章封面id
	Tags     enum.Array `json:"tags"`                                          // 文章标签
}

func ArticleCreateService(cr ArticleAddRequest, claims *jwts.MyClaims) (err error) {

	UserID := claims.UserID
	UserNickName := claims.Username

	// 处理content  原始 Markdown → 转 HTML → 检查并移除危险的<script>标签 → 转回 Markdown → 替换原始内容
	unsafe := blackfriday.MarkdownCommon([]byte(cr.Content))
	// 是不是有script标签
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	//fmt.Println(doc.Text())  全是字了
	nodes := doc.Find("script").Nodes
	if len(nodes) > 0 {
		// 有script标签
		doc.Find("script").Remove()
		converter := md.NewConverter("", true, nil)
		html, _ := doc.Html()
		markdown, _ := converter.ConvertString(html)
		cr.Content = markdown
	}

	if cr.Abstract == "" {
		// 汉字的截取不一样
		abs := []rune(doc.Text())
		// 将content转为html，并且过滤xss，以及获取中文内容
		if len(abs) > 100 {
			cr.Abstract = string(abs[:100])
		} else {
			cr.Abstract = string(abs)
		}
	}
	if cr.BannerID == 0 {
		//说明没传  随机在数据库中选择一个
		var BannerModels []models.BannerModel
		global.DB.Model(&models.BannerModel{}).Find(&BannerModels)
		if len(BannerModels) == 0 {
			//数据库一个照片都没有
			return errors.New("数据库未有图片")
		}

		// 生成 0 到 lens 之间的随机整数
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNum := r.Intn(len(BannerModels)) // Intn(n) 生成 [0, n) 范围内的整数
		cr.BannerID = BannerModels[randomNum].ID
	}
	var bannerUrl string
	global.DB.Model(models.BannerModel{}).Where("id = ?", cr.BannerID).Select("path").Scan(&bannerUrl)

	// 查用户头像
	var avatar string
	err = global.DB.Model(models.UserModel{}).Where("id = ?", UserID).Select("avatar").Scan(&avatar).Error
	if err != nil {
		return fmt.Errorf("用户不存在 %s", err.Error())
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	article := models.ArticleModel{
		CreatedAt:    now,
		UpdatedAt:    now,
		Title:        cr.Title,
		Keyword:      cr.Title, //Keyword 必然是精确匹配
		Abstract:     cr.Abstract,
		Content:      cr.Content,
		UserID:       UserID,
		UserNickName: UserNickName,
		UserAvatar:   avatar,
		Category:     cr.Category,
		Source:       cr.Source,
		Link:         cr.Link,
		BannerID:     cr.BannerID,
		BannerUrl:    bannerUrl,
		Tags:         cr.Tags,
	}
	if article.ISExistData() {
		//已经存在
		return errors.New("文章已经存在")
	}

	err = article.Create()
	if err != nil {
		return fmt.Errorf("创建文章失败 %s", err.Error())
	}
	return nil
}
