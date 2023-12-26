package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{
		// 慢查询logger
		// Logger:
	})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	return db
}
