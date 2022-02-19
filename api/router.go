package api

import (
	"coinche/usecases"
	"coinche/utilities"

	"github.com/gin-gonic/gin"
)

func SetupRouter(gameUsecases *usecases.GameUsecases) (*gin.Engine, *Hub) {
	gameAPIs := &GameAPIs{Usecases: gameUsecases}

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	utilities.PanicIfErr(err)

	hub := NewHub()
	go hub.run()

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.GET("/games/:id/join", func(c *gin.Context) {
		gameAPIs.JoinGame(c, hub)
	})

	return router, hub
}
