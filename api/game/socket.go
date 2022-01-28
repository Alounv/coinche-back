package gameapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GameSocketHandler(context *gin.Context) {
	HTTPGameSocketHandler(context.Writer, context.Request)
}

func HTTPGameSocketHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		t, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = conn.WriteMessage(t, message)
		if err != nil {
			break
		}
	}
}
