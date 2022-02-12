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
		context.JSON(http.StatusBadRequest, gin.H{"error": "WRONG ID FORMAT"})
		return
	}
	game, err := gameAPIs.Usecases.GetGame(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "GAME NOT FOUND"})
		return
	}

	context.JSON(http.StatusOK, game)
}
