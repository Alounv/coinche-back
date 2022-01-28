package gameapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetGame is exported for other packages such as repository and api
func (gameAPIs *GameAPIs) GetGame(context *gin.Context) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
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
