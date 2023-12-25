package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{
		// 慢查询logger
		// Logger:
	})
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
	return db
}
