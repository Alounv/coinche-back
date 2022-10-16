package api

import (
	"coinche/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(gameUsecases *usecases.GameUsecases, origins []string) (*gin.Engine, *Hub) {
	gameAPIs := &GameAPIs{Usecases: gameUsecases}

	router := gin.Default()

	config := cors.DefaultConfig()
	if len(origins) >= 1 {
		config.AllowOrigins = origins
	} else {
		config.AllowAllOrigins = true
	}
	config.AllowMethods = []string{"PUT", "PATCH", "GET", "DELETE", "POST"}

	router.Use(cors.New(config))

	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	hub := NewHub(gameUsecases)
	go hub.run()

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.DELETE("/games/:id/delete", gameAPIs.deleteGame)
	router.PATCH("/games/:id/archive", gameAPIs.archiveGame)
	router.PUT("/games/:id/leave", gameAPIs.leaveGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.GET("/games/:id/join", func(c *gin.Context) {
		gameAPIs.JoinGame(c, &hub)
	})

	return router, &hub
}
