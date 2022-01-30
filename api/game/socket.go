package gameapi

import (
	"encoding/json"
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
		panic(err)
	}

	err = SendMessage(conn, "connection established")
	if err != nil {
		panic(err)
	}

	for {
		message, err := ReceiveMessage(conn)
		if err != nil {
			break
		}
		err = SendMessage(conn, message)
		if err != nil {
			break
		}
	}
}

func SendMessage(connection *websocket.Conn, msg string) error {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = connection.WriteMessage(websocket.BinaryMessage, message)

	return err
}

func ReceiveMessage(connection *websocket.Conn) (string, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return "", err
	}

	var reply string
	err = json.Unmarshal(message, &reply)

	return reply, err
}
