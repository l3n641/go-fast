package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func main() {
	initConfig()
	initFileLog()
	router := gin.Default()
	debug := viper.GetString("app.debug")
	if debug != "" {
		gin.SetMode(gin.DebugMode)
	}

	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, router)
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