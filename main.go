package main

import (
	"github.com/tofustream/simple-auth-service/config"
	"github.com/tofustream/simple-auth-service/routes"
)

func main() {
	// DB 接続
	config.ConnectDB()
	r := routes.SetupRouter()
	r.Run(":8000")
}
