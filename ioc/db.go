package ioc

import (
	"github.com/jasonzhao47/cuddle/configs"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(configs.Config.DB.DSN))

	if err != nil {
		panic(err)
	}
	// should init tables here at dao layer
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate
	return db
}
