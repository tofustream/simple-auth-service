package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tofustream/simple-auth-service/config"
	"github.com/tofustream/simple-auth-service/routes"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	config.LoadCORSConfig()
	log.Println("Allowed Origins:", config.AllowOrigins)

	config.ConnectDB()

	r := routes.SetupRouter()

	// 環境変数 PORT の値を取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // デフォルトポートを設定
	}

	// サーバーを指定したポートで起動
	r.Run(":" + port)
}
