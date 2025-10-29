package image_service

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"Blog_server/models"
	"Blog_server/models/enum"
	"Blog_server/plugins/qiniu"
	"Blog_server/utils/image"
	"Blog_server/utils/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path"
	"strings"
)

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

type ImageListResponse struct {
	Filename  string  `json:"filename"`
	FilePath  string  `json:"filepath"`
	Size      float64 `json:"size"`
	IsSuccess bool    `json:"is_success"`
	Msg       string  `json:"msg"`
}

func UploadService(c *gin.Context, FileHeader *multipart.FileHeader) (ImageListResponse, error) {
	//首先 先判断白名单
	filename := FileHeader.Filename
	splitList := strings.Split(filename, ".")
	suffix := strings.ToLower(splitList[(len(splitList) - 1)])

	ok := image.InList(suffix, WhiteImageList)
	if !ok {
		return ImageListResponse{
			Filename:  FileHeader.Filename,
			Size:      float64(FileHeader.Size) / float64(1024*1024),
			IsSuccess: false,
			Msg:       "图片格式错误 上传失败",
		}, nil
	}
	//如果超出 预设的2MB  上传失败  如果没超 那就上传成功
	size := float64(FileHeader.Size) / float64(1024*1024) //size 单位 MB
	if size <= float64(global.Config.Upload.Size) {
		//可以上传  设置上传本地还是 七牛
		imageType := enum.LocationType //默认
		filepath := path.Join("uploads", FileHeader.Filename)

		// md5  去数据库查看hash 存不存在  哈希算法
		hashString, ByteData, err := md5.MD5_Hash(FileHeader)
		if err != nil { //转换失败
			res.FailWithErr(c, err)
			return ImageListResponse{}, err
		}

		err = global.DB.Take(&models.BannerModel{}, "Hash = ?", hashString).Error
		if err == nil {
			//说明已经找到了  上传重复
			return ImageListResponse{
				Filename:  FileHeader.Filename,
				FilePath:  filepath,
				Size:      size,
				IsSuccess: false,
				Msg:       "图片已存在 上传失败",
			}, nil
		}

		msg := fmt.Sprintf("该图片共%.2fMB， 上传成功", size)

		if global.Config.QiNiu.Enable {
			//上传 QiNiu
			imageType = enum.QiNiuType
			filepath, err = qiniu.UploadImage(ByteData, filename, global.Config.QiNiu.Prefix)
			if err != nil {
				return ImageListResponse{}, err
			}
			msg = "上传七牛云成功"

		} else {
			err = c.SaveUploadedFile(FileHeader, filepath)
			if err != nil {
				return ImageListResponse{}, err
			}
		}

		//存进数据库
		global.DB.Create(&models.BannerModel{
			Path:      filepath,
			Hash:      hashString,
			Name:      filename,
			ImageType: imageType,
		})

		return ImageListResponse{
			Filename:  FileHeader.Filename,
			FilePath:  filepath,
			Size:      size,
			IsSuccess: true,
			Msg:       msg,
		}, nil

	} else {
		//说明 文件大小过大
		return ImageListResponse{
			Filename:  FileHeader.Filename,
			Size:      size,
			IsSuccess: false,
			Msg:       fmt.Sprintf("该图片共%.2fMB，超出预设的%dMB,上传失败", size, global.Config.Upload.Size),
		}, nil

	}

}
