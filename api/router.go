package api

import (
	"coinche/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(gameUsecases *usecases.GameUsecases) *gin.Engine {
	gameAPIs := &GameAPIs{Usecases: gameUsecases}

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	hub := newHub()
	go hub.run()

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.GET("/games/:id/join", gameAPIs.JoinGame)
	router.GET("/games/:id/join2", func(c *gin.Context) {
		gameAPIs.JoinGame2(c, hub)
	})

	return router
}
