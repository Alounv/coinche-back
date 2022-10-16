package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		connectionOrigin := r.Header.Get("Origin")
		fmt.Println(connectionOrigin)
		if connectionOrigin == "" {
			return true
		}

		authorizedOrigin := os.Getenv("AUTHORIZED_ORIGIN")
		return connectionOrigin == authorizedOrigin
	},
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
	if err != nil {
		fmt.Println("Error upgrading socket with new player: ", err)
		return
	}

	PlayerSocketHandler(connection, gameID, playerName, hub)
}
