package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) deleteGame(context *gin.Context) {
	stringID := context.Query("id")
	gameID, err := strconv.Atoi(stringID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	err = gameAPIs.Usecases.DeleteGame(gameID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gameID)
}
