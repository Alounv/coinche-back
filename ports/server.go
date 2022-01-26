package ports

import (
	"coinche/app"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Store app.GameServiceType
}

func SetupRouter(store app.GameServiceType) *gin.Engine {
	server := &Server{store}

	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/games/:id", server.GetGame)
	router.POST("/games/:name", server.CreateGame)
	router.GET("/games/all", server.ListGames)

	return router
}
