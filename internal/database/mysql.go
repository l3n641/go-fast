package database

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

type MysqlDatabase struct {
	instance *gorm.DB
	once     sync.Once
}

func NewMysqlDatabase() *MysqlDatabase {
	return &MysqlDatabase{}
}

func (d *MysqlDatabase) GetInstance() *gorm.DB {
	d.once.Do(func() {
		var defaultLogger logger.Interface
		dsn := viper.GetString("database_mysql.dsn")
		debugMode := viper.GetString("database_mysql.debug")
		if debugMode != "" {
			defaultLogger = logger.Default.LogMode(logger.Info)
		}

		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: defaultLogger,
		})
		if err != nil {
			panic(err)
		}
		d.instance = conn
	})

	return d.instance
}
