package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/spf13/viper"
	"go-fast/api/middleware"
	"go-fast/internal/database"
	"go-fast/open.go/webscockt"
	_ "go-fast/websocketApi"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	initConfig()
	initFileLog()
	initMysql()
	router := gin.Default()
	debug := viper.GetString("app.debug")
	if debug != "" {
		gin.SetMode(gin.DebugMode)
	}
	router.Use(middleware.Cors)

	m := initWebsocket()
	// WebSocket 路由
	router.GET("/ws", func(c *gin.Context) {
		if c.GetHeader("Upgrade") != "websocket" {
			return
		}
		err := m.HandleRequest(c.Writer, c.Request)
		if err != nil {
			// 处理错误
		}
	})

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

func initMysql() {
	db := database.NewMysqlDatabase().GetInstance()
	db.AutoMigrate()
}

func initRedis() {
	database.NewRedisClient()
}

func initWebsocket() *webscockt.WS {
	config := &melody.Config{
		WriteWait:         10 * time.Second,
		PongWait:          20 * time.Second,
		PingPeriod:        5 * time.Second,
		MaxMessageSize:    512,
		MessageBufferSize: 256,
	}
	return webscockt.NewWebsocket(config).GetWebsocket()

}
