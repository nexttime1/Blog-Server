package image_api

import (
	"Blog_server/common"
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/utils/image"
	"Blog_server/utils/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

type ImageApi struct {
}

type ImageListResponse struct {
	Filename  string  `json:"filename"`
	Size      float64 `json:"size"`
	IsSuccess bool    `json:"is_success"`
	Msg       string  `json:"msg"`
}

// ImageListView 接受参数
type ImageListView struct {
	common.PageInfo
	Path string `form:"path"` // 图片路径
	Hash string `form:"hash"` // 图片的hash值，用于判断重复图片
	Name string `form:"name"` // 图片名称
}

// WhiteImageList 白名单
var WhiteImageList = []string{
	"jpg",
	"png",
	"jpeg",
	"ico",
	"tiff",
	"gif",
	"svg",
	"webp",
}

func (ImageApi) ImageUploadView(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	FileHeaderList, ok := form.File["image"]
	if !ok {
		res.FailWithErr(c, errors.New("上传的文件图片不存在"))
		return
	}
	// 记录n个照片上传
	var response []ImageListResponse

	for _, FileHeader := range FileHeaderList {
		//首先 先判断白名单
		filename := FileHeader.Filename
		splitList := strings.Split(filename, ".")
		suffix := strings.ToLower(splitList[(len(splitList) - 1)])

		ok = image.InList(suffix, WhiteImageList)
		if !ok {
			response = append(response, ImageListResponse{
				Filename:  FileHeader.Filename,
				Size:      float64(FileHeader.Size) / float64(1024*1024),
				IsSuccess: false,
				Msg:       "图片格式错误 上传失败",
			})
			continue
		}

		//如果超出 预设的2MB  上传失败  如果没超 那就上传成功
		size := float64(FileHeader.Size) / float64(1024*1024) //size 单位 MB
		if size <= float64(global.Config.Upload.Size) {
			//可以上传
			filepath := path.Join("uploads", FileHeader.Filename)
			// md5  去数据库查看hash 存不存在
			hashString, err := md5.MD5_Hash(c, FileHeader)
			if err != nil {
				res.FailWithErr(c, err)
				return
			}
			err = global.DB.Take(&models.BannerModel{}, "Hash = ?", hashString).Error
			if err == nil {
				//说明已经找到了  上传重复
				response = append(response, ImageListResponse{
					Filename:  FileHeader.Filename,
					Size:      float64(FileHeader.Size) / float64(1024*1024), //size 单位 MB
					IsSuccess: false,
					Msg:       "图片已存在上传失败",
				})
				continue
			}
			global.DB.Create(&models.BannerModel{
				Path: filepath,
				Hash: hashString,
				Name: filename,
			})

			err = c.SaveUploadedFile(FileHeader, filepath)
			if err != nil {
				res.FailWithErr(c, err)
				return
			}

			response = append(response, ImageListResponse{
				Filename:  FileHeader.Filename,
				Size:      size,
				IsSuccess: true,
				Msg:       fmt.Sprintf("该图片共%.2fMB， 上传成功", size),
			})
		} else {
			response = append(response, ImageListResponse{
				Filename:  FileHeader.Filename,
				Size:      size,
				IsSuccess: false,
				Msg:       fmt.Sprintf("该图片共%.2fMB，超出预设的%dMB,上传失败", size, global.Config.Upload.Size),
			})
		}

	}
	res.OkWithData(c, response)
}

// ImageInfoView 查看图片
func (ImageApi) ImageInfoView(c *gin.Context) {
	var cr ImageListView
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	list, count, err := common.ListQuery(models.BannerModel{
		Path: cr.Path,
		Hash: cr.Hash,
		Name: cr.Name,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Debug:        true,
		DefaultOrder: "created_at DESC", //写死了  默认降序排序  前端修改的话
	})
	if err != nil {
		res.FailWithErr(c, err)
		return
	}
	res.OkWithList(c, list, count)

}
