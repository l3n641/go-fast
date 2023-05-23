package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-fast/api/middleware"
	"net/http"
)

var serverCmd = &cobra.Command{
	Use:     "http-server",
	Short:   "启动http服务",
	Example: "go-fast http-server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	router := gin.Default()
	debug := viper.GetString("app.debug")
	if debug != "" {
		gin.SetMode(gin.DebugMode)
	}
	router.Use(middleware.Cors)

	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, router)
}
