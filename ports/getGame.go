package ports

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetGame(context *gin.Context) {
	stringId := context.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		panic(err)
	}
	game := server.Store.GetGame(id)

	if game.Name != "" {
		context.JSON(http.StatusOK, game)
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
	}
}
