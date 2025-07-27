package md5

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
)

func MD5(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func MD5_Hash(FileHeader *multipart.FileHeader) (string, []byte, error) {
	fileObj, err := FileHeader.Open()
	if err != nil {
		return "", nil, errors.New("打开 FileHeader 失败")
	}
	ByteData, err := io.ReadAll(fileObj)
	if err != nil {
		return "", nil, errors.New("ReadAll 读取失败")
	}
	return MD5(ByteData), ByteData, nil
}
