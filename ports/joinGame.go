package ports

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) JoinGame(context *gin.Context) {
	stringId := context.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		panic(err)
	}

	playerName := context.Query("playerName")

	err = server.Store.JoinGame(id, playerName)

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
			context.Status(http.StatusInternalServerError)
		}
	}
}
