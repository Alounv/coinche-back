package ports

import (
	"coinche/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Store app.GameServiceType
}

func SetupRouter(store app.GameServiceType) *gin.Engine {
	server := &Server{store}

	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/games/:id", server.GetAGame)
	router.POST("/games/:name", server.CreateAGame)
	router.GET("/games/all", server.GetAllGames)

	return router
}

func (server *Server) GetAGame(context *gin.Context) {
	stringId := context.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		panic(err)
	}
	game := server.Store.GetAGame(id)

	if game.Name != "" {
		context.JSON(http.StatusOK, game)
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
	}
}

func (server *Server) CreateAGame(context *gin.Context) {
	name := context.Param("name")

	server.Store.CreateAGame(name)
	context.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

func (server *Server) GetAllGames(context *gin.Context) {
	games := server.Store.GetAllGames()

	context.JSON(http.StatusOK, games)
}
