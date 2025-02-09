package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tofustream/simple-auth-service/config"
	"github.com/tofustream/simple-auth-service/handlers"
)

const (
	registerRoute = "/register"
	loginRoute    = "/login"
	refreshRoute  = "/refresh"
	logoutRoute   = "/logout"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 設定を適用
	r.Use(cors.New(config.GetCORSConfig()))

	r.POST(registerRoute, handlers.Register)
	r.POST(loginRoute, handlers.Login)
	r.POST(refreshRoute, handlers.RefreshToken)
	r.POST(logoutRoute, handlers.Logout)

	return r
}
