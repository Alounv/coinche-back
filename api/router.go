package api

import (
	gameapi "coinche/api/game"
	"coinche/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(gameService domain.GameServiceType) *gin.Engine {
	gameAPIs := &gameapi.GameAPIs{GameService: gameService}

	router := gin.Default()
	err := router.SetTrustedProxies([]string{"192.168.1.2"})
	if err != nil {
		panic(err)
	}

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.POST("/games/:id/join", gameAPIs.JoinGame)

	return router
}
