package gameapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) GetGame(context *gin.Context) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
	if err != nil { // CONTEXT.JSON()...
		fmt.Println("ID format is wrong", err)
	}
	game, err := gameAPIs.Usecases.GetGame(id)
	if err != nil {
		fmt.Println("Game not found", err)
	}

	if game.Name != "" {
		context.JSON(http.StatusOK, game)
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "GAME NOT FOUND"})
	}
}
