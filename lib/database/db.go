/*
*

	@author:
	@date : 2025/5/15
*/
package database

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
)

var db *gorm.DB

// Register 初始化数据库客户端
func Register() {

	//数据库注册
	var err error
	driver := viper.GetString("db.driver")
	source := viper.GetString("db.dsn")
	db, err = gorm.Open(driver, source)
	if err != nil {
		log.Panic(err)
	}
	db.LogMode(true)
	////空闲最大连接数
	//db.DB().SetMaxIdleConns(10)
	////设置打开数据库连接的最大数量。
	//db.DB().SetMaxOpenConns(100)
}

func GetDB() *gorm.DB {
	return db
}
