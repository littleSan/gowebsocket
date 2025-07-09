/*
@author: little
@date : 2023/10/8
*/
package ossCli

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"log"
)

var ossClient = &oss.Client{}

func NewOssCli() {
	fmt.Println("OSS Go SDK Version: ", oss.Version)

	var err error
	ossClient, err = oss.New(viper.GetString("aliOss.Endpoint"), viper.GetString("aliOss.Id"), viper.GetString("aliOss.Secret"))
	if err != nil {
		log.Println("初始化aliOSS 客户端失败", err)
		return
	}
	fmt.Println("初始化aliOSS 成功")
}
func GetOssCli() *oss.Client {
	return ossClient
}
