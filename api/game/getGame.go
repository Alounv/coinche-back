package gameapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) GetGame(context *gin.Context) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		panic(err)
	}
	game := gameAPIs.GameService.GetGame(id)

	if game.Name != "" {
		context.JSON(http.StatusOK, game)
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
	}
}
