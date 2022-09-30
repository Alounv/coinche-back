package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) leaveGame(context *gin.Context) {
	stringID := context.Param("id")
	gameID, err := strconv.Atoi(stringID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	playerName := context.Query("playerName")

	err = gameAPIs.Usecases.LeaveGame(gameID, playerName)
	fmt.Println(err)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusAccepted, gameID)
}
