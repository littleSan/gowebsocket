/*
*

	@author:
	@date : 2025/6/3
*/
package fileControl

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/controllers"
	"github.com/link1st/gowebsocket/v2/helper/validate"
	"github.com/link1st/gowebsocket/v2/lib/ossCli"
	"github.com/spf13/viper"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// 文件上传下载
func UploadPic(ctx *gin.Context) {
	// 支持多文件上传
	multiForm, err := ctx.MultipartForm()
	if err != nil {
		fmt.Println("解析上传文件参数出错")
		controllers.Response(ctx, common.ParameterIllegal, err.Error(), nil)
		return
	}
	//获取文件
	fileSlice := []string{}
	f := multiForm.File["file"]
	if !validate.FilesValidate(f, validate.IMGAllow, int64(validate.FileMaxSize)) {
		fmt.Println("文件格式错误")
		controllers.Response(ctx, common.FileFormatError, err.Error(), nil)
		return
	}
	data := make(map[string]interface{})
	fileSize := len(f)

	for _, multiFile := range f {
		fmt.Println("上传文件名称", multiFile.Filename)
		reset, err := UploadFile("im/picture", multiFile)
		if err != nil {
			continue
		}
		fileSlice = append(fileSlice, reset)
		if fileSize == 1 {
			data["picUrl"] = fileSlice
			controllers.Response(ctx, common.OK, "", data)
			return
		}

	}
	controllers.Response(ctx, common.OK, "", data)
	return
}

func UploadFile(path string, file *multipart.FileHeader) (rest string, err error) {
	dateStr := time.Now().Format("2006-01-02")
	folderName := strings.ReplaceAll(dateStr, "-", "")
	var uploadFileName = path + "/" + folderName + "/"
	isExist, err := ossCli.GetOssCli().IsBucketExist(viper.GetString("aliOss.Bucket"))
	if err != nil {
		fmt.Println("Error:", err)
		return "", errors.New("bucket err")
	}
	fmt.Println("IsBucketExist result : ", isExist)
	//文件上传，文件上传有简单上传，追加上传，断点续传上传，分片上传
	if !isExist {
		return "", errors.New("bucket not exist")
	}
	//注意此处不要写错，写错的话，err让然是nil，我们应该需要先判断一下是否存在
	bucket, _ := ossCli.GetOssCli().Bucket(viper.GetString("aliOss.Bucket"))
	source, err := file.Open()
	if err != nil {
		fmt.Println("open file err", err)
		return "", errors.New("open file err")
	}
	defer source.Close()
	newName := getHash(file.Filename + time.Now().String())
	uploadFileName = uploadFileName + newName + filepath.Ext(file.Filename)
	err = bucket.PutObject(uploadFileName, source)
	if err != nil {
		fmt.Println("path", uploadFileName)
		return "", errors.New("upload file err")
	}
	return getPath(uploadFileName), err
}

func UploadUrl(path string, url string) (rest string, err error) {
	dateStr := time.Now().Format("2006-01-02")
	folderName := strings.ReplaceAll(dateStr, "-", "")
	var uploadFileName = path + "/" + folderName + "/"
	isExist, err := ossCli.GetOssCli().IsBucketExist(viper.GetString("aliOss.Bucket"))
	if err != nil {
		fmt.Println("Error:", err)
		return "", errors.New("bucket err")
	}
	fmt.Println("IsBucketExist result : ", isExist)
	//文件上传，文件上传有简单上传，追加上传，断点续传上传，分片上传
	if !isExist {
		return "", errors.New("bucket not exist")
	}
	//注意此处不要写错，写错的话，err让然是nil，我们应该需要先判断一下是否存在
	bucket, _ := ossCli.GetOssCli().Bucket(viper.GetString("aliOss.Bucket"))
	//获取URL文件
	resp, err := http.Get(url) // 从URL获取图片内容
	if err != nil {
		fmt.Println("Failed to retrieve image from URL:", err)
		return
	}
	defer resp.Body.Close()
	nowName := getHash(url + time.Now().String())
	uploadFileName = uploadFileName + nowName + filepath.Ext(url)
	bucket.PutObject(uploadFileName, resp.Body)
	return getPath(uploadFileName), nil

}

func getHash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func getPath(path string) string {
	return "https://" + viper.GetString("aliOss.Bucket") + "." + viper.GetString("aliOss.Endpoint") + "/" + path
}
