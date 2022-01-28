package gameApi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) GetGame(context *gin.Context) {
	stringId := context.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		panic(err)
	}
	game := gameAPIs.Store.GetGame(id)

	if game.Name != "" {
		context.JSON(http.StatusOK, game)
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
	}
}
