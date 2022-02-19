package api

import (
	"coinche/utilities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (gameAPIs *GameAPIs) JoinGame(context *gin.Context, hub *Hub) {
	stringID := context.Param("id")
	gameID, err := strconv.Atoi(stringID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	playerName := context.Query("playerName")

	connection, err := wsupgrader.Upgrade(context.Writer, context.Request, nil)
	utilities.PanicIfErr(err)
	PlayerSocketHandler(connection, gameAPIs.Usecases, gameID, playerName, hub)
}
