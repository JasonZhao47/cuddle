package ioc

import (
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var c Config
	c = Config{
		DSN: "root:root@tcp(127.0.0.1:3306)/cuddle",
	}
	err := viper.UnmarshalKey("data.database", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败%s", err.Error()))
	}

	db, err := gorm.Open(mysql.Open(c.DSN))
	if err != nil {
		panic(err)
	}
	// should init tables on dao layer
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate
	return db
}
