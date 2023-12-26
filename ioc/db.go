package ioc

import (
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/pkg/gormx"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
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
	err = db.Use(prometheus.New(
		prometheus.Config{
			DBName:          "cuddle",
			RefreshInterval: 15,
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.MySQL{
					// 配置有多少正在运行的线程
					VariableNames: []string{"Threads_running"},
				},
			},
		},
	))
	if err != nil {
		panic(err)
	}
	cb := gormx.NewCallback(prometheus2.SummaryOpts{
		Namespace: "jason_zhao",
		Subsystem: "cuddle",
		Name:      "gorm_db_",
		Help:      "统计 GORM 的数据库查询",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		}})
	err = db.Use(cb)
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate
	return db
}
