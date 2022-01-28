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

	err = gameAPIs.Usecases.JoinGame(id, playerName)

	if err == nil {
		context.Status(http.StatusAccepted)
		return
	}

	switch err.Error() {
	case "GAME NOT FOUND":
		{
			context.Status(http.StatusNotFound)
		}
	case "GAME IS FULL":
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
