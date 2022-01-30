package gameapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) JoinGame(context *gin.Context) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	playerName := context.Query("playerName")

	_, err = gameAPIs.Usecases.JoinGame(id, playerName)

	if err == nil {
		context.Status(http.StatusAccepted)
		return
	}

	switch err.Error() {
	case "GAME NOT FOUND":
		{
			context.Status(http.StatusNotFound)
			return
		}
	case "GAME IS FULL":
		{
			context.Status(http.StatusForbidden)
			return
		}
	default:
		{
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
