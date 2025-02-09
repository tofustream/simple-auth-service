package config

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

// AllowOrigins は CORS で許可するオリジンのリスト
var AllowOrigins []string

// LoadCORSConfig は CORS 設定をロードする
func LoadCORSConfig() {
	_ = godotenv.Load()

	// 環境変数から許可するオリジンを取得（カンマ区切りでリスト化）
	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins != "" {
		AllowOrigins = strings.Split(origins, ",")
	} else {
		AllowOrigins = []string{"*"} // デフォルトは全て許可
	}
}

// GetCORSConfig は Gin の CORS 設定を返す
func GetCORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
}
