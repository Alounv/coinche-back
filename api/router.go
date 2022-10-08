package api

import (
	"coinche/usecases"
	"coinche/utilities"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(gameUsecases *usecases.GameUsecases, origins []string) (*gin.Engine, *Hub) {
	gameAPIs := &GameAPIs{Usecases: gameUsecases}

	router := gin.Default()

	config := cors.DefaultConfig()
	//fmt.Println("origins", origins) // we should understand why it's not working
	//config.AllowOrigins = []string{"http://127.0.0.1:5173"}
	//config.AllowOrigins = origins
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	err := router.SetTrustedProxies(nil)
	utilities.PanicIfErr(err)

	hub := NewHub(gameUsecases)
	go hub.run()

	router.GET("/games/:id", gameAPIs.GetGame)
	router.POST("/games/create", gameAPIs.CreateGame)
	router.DELETE("/games/:id/delete", gameAPIs.deleteGame)
	router.PUT("/games/:id/leave", gameAPIs.leaveGame)
	router.GET("/games/all", gameAPIs.ListGames)
	router.GET("/games/:id/join", func(c *gin.Context) {
		gameAPIs.JoinGame(c, &hub)
	})

	return router, &hub
}
