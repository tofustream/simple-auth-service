package main

import (
	"log"

	"github.com/tofustream/simple-auth-service/config"
	"github.com/tofustream/simple-auth-service/routes"
)

func main() {
	config.LoadCORSConfig()
	log.Println("Allowed Origins:", config.AllowOrigins)

	config.ConnectDB()

	r := routes.SetupRouter()
	r.Run(":8000")
}
