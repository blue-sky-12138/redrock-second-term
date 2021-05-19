package MySQL

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DNS = "bluesky:135246Cjw@tcp(rm-bp14fk5x3q4byb6a2so.mysql.rds.aliyuncs.com:3306)" +
	"/oauth2?charset=utf8mb4&parseTime=True&loc=Local"

var (
	DB *gorm.DB
)

func MySQLInit() error {
	db, err := gorm.Open(mysql.Open(DNS))
	if err != nil {
		return err
	}
	DB = db

	err = AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

func MySQLDebug() error {
	db, err := gorm.Open(mysql.Open(DNS), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return err
	}
	DB = db.Debug()

	err = AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

//自动迁移
func AutoMigrate() error {
	err := DB.AutoMigrate(
	//&ClientInfo{},
	)
	return err
}
