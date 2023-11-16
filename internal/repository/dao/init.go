package dao

import "gorm.io/gorm"

// initialize tables
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Article{})
}
