package gameapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (gameAPIs *GameAPIs) JoinGame(context *gin.Context) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		panic(err)
	}

	playerName := context.Query("playerName")

	err = gameAPIs.GameService.JoinGame(id, playerName)

	if err == nil {
		context.Status(http.StatusAccepted)
		return
	}

	switch err.Error() {
	case "Game not found":
		{
			context.Status(http.StatusNotFound)
		}
	case "Game is full":
		{
			context.Status(http.StatusForbidden)
		}
	default:
		{
			fmt.Print(err)
			context.Status(http.StatusInternalServerError)
		}
	}
}