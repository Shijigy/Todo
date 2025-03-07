package utils

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func UploadImageToQiNiu(r *http.Request) (string, error) {
	// 获取七牛云的上传凭证和上传目录
	accessKey := "D3BspotCT7UQRJx0q8GaznFvHTPJ-AVQC_IFPmjv"
	secretKey := "t0Ys1TnexpaiaOQLVMeeubyBH0HOvtkLq1SlHXVp"
	bucket := "todo22"

	// 获取上传文件
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "upload_")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	// 将multipart.File的内容写入到临时文件
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return "", err
	}

	// 获取文件的大小
	fileStat, err := tmpFile.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileStat.Size() // 获取文件大小，返回 int64

	// 七牛云上传凭证
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	// 设置上传配置
	cfg := storage.Config{
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)

	// 获取 Gin 请求的 context
	ctx := r.Context()

	// 上传文件到七牛云
	fileKey := fmt.Sprintf("post-images/%s", generateFileName()) // 生成文件名
	ret := storage.PutRet{}
	err = formUploader.Put(ctx, &ret, upToken, fileKey, tmpFile, fileSize, nil)
	if err != nil {
		return "", err
	}

	// 返回文件的 URL
	fileURL := fmt.Sprintf("http://%s/%s", "ssjwo2ece.hn-bkt.clouddn.com", ret.Key)
	return fileURL, nil
}

// generateFileName 生成一个唯一的文件名
func generateFileName() string {
	// 获取当前时间戳
	timestamp := time.Now().UnixNano()

	// 生成一个随机字符串
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	var randomString strings.Builder
	for i := 0; i < length; i++ {
		randomString.WriteByte(charset[rand.Intn(len(charset))])
	}

	// 组合时间戳和随机字符串作为文件名
	fileName := fmt.Sprintf("%d_%s", timestamp, randomString.String())

	return fileName
}
