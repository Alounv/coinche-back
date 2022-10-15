package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		connectionOrigin := r.Header.Get("Origin")
		if connectionOrigin == "" {
			return true
		}

		//authorizedOrigin := os.Getenv("AUTHORIZED_ORIGIN")
		// fmt.Println("------", authorizedOrigin, connectionOrigin) // we should understand why it's not working
		// return connectionOrigin == "http://127.0.0.1:5173" || connectionOrigin == "http://localhost:5000"
		return true
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
