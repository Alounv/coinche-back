package api

import (
	gameApi "coinche/api/game"
	"coinche/domain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(gameService domain.GameServiceType) *gin.Engine {
	gameAPIs := &gameApi.GameAPIs{Store: gameService}

	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.POST("/games/:id/join", gameAPIs.JoinGame)

	return router
}
