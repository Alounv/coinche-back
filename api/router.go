package api

import (
	gameapi "coinche/api/game"
	"coinche/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(gameUsecases *usecases.GameUsecases) *gin.Engine {
	gameAPIs := &gameapi.GameAPIs{Usecases: gameUsecases}

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.GET("/games/:id/join", gameAPIs.JoinGame)

	return router
}
