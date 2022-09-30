package api

import (
	"coinche/usecases"
	"coinche/utilities"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(gameUsecases *usecases.GameUsecases) (*gin.Engine, *Hub) {
	gameAPIs := &GameAPIs{Usecases: gameUsecases}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:5173"}
	router.Use(cors.New(config))

	err := router.SetTrustedProxies(nil)
	utilities.PanicIfErr(err)

	hub := NewHub(gameUsecases)
	go hub.run()

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.GET("/games/:id/join", func(c *gin.Context) {
		gameAPIs.JoinGame(c, &hub)
	})

	return router, &hub
}
