package image_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
)

type ImageApi struct {
}

type ImageListResponse struct {
	Filename  string  `json:"filename"`
	Size      float64 `json:"size"`
	IsSuccess bool    `json:"is_success"`
	Msg       string  `json:"msg"`
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
		//如果超出 预设的2MB  上传失败  如果没超 那就上传成功
		size := float64(FileHeader.Size) / float64(1024*1024) //size 单位 MB
		if size <= float64(global.Config.Upload.Size) {
			//可以上传
			filepath := path.Join("uploads", FileHeader.Filename)
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
