package gameapi

import (
	"coinche/domain"
	"coinche/usecases"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (gameAPIs *GameAPIs) JoinGame(context *gin.Context) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	playerName := context.Query("playerName")

	HTTPGameSocketHandler(context.Writer, context.Request, gameAPIs.Usecases, id, playerName)
}

func HTTPGameSocketHandler(
	writer http.ResponseWriter,
	request *http.Request,
	usecases usecases.GameUsecasesInterface,
	id int,
	playerName string,
) {
	connection, err := wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}

	game, err := usecases.JoinGame(id, playerName)
	if err != nil {
		err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err))
		if err != nil {
			panic(err)
		}
		connection.Close()
	}

	err = sendGame(connection, game)
	if err != nil {
		panic(err)
	}

	for {
		message, err := ReceiveMessage(connection)
		if err != nil {
			break
		}
		if message == "leave" {
			err = usecases.LeaveGame(id, playerName)
			if err != nil {
				fmt.Println("Could not leave this game: ", err)
				break
			}
			err = SendMessage(connection, "Has left the game")
			connection.Close()
		} else {
			err = SendMessage(connection, message)
		}
		if err != nil {
			break
		}
	}
}

func sendGame(connection *websocket.Conn, game domain.Game) error {
	message, err := json.Marshal(game)
	if err != nil {
		return err
	}

	err = send(connection, message)
	return err
}

func SendMessage(connection *websocket.Conn, msg string) error {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = send(connection, message)
	return err
}

func send(connection *websocket.Conn, message []byte) error {
	err := connection.WriteMessage(websocket.BinaryMessage, message)
	return err
}

func ReceiveGame(connection *websocket.Conn) (domain.Game, error) {
	var game domain.Game
	message, err := receive(connection)
	if err != nil {
		return game, err
	}

	err = json.Unmarshal(message, &game)
	return game, err
}

func ReceiveMessage(connection *websocket.Conn) (string, error) {
	var reply string
	message, err := receive(connection)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(message, &reply)
	return reply, err
}

func receive(connection *websocket.Conn) ([]byte, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, err
}
