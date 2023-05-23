package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-fast/cmd"
	"go-fast/internal/database"
	"io"
	"os"
)

func main() {
	initConfig()
	initFileLog()
	initMysql()

	cmd.Execute()

}

func initConfig() {
	viper.SetConfigName("configs/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

// 初始化日志
func initFileLog() {
	gin.DisableConsoleColor()
	logFile := viper.GetString("app.logFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initMysql() {
	db := database.NewMysqlDatabase().GetInstance()
	db.AutoMigrate()
}

func initRedis() {
	database.NewRedisClient()
}
