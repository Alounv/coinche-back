package api

import (
	"coinche/domain"
	"coinche/usecases"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (gameAPIs *GameAPIs) JoinGame2(context *gin.Context, hub *hub) {
	stringID := context.Param("id")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	playerName := context.Query("playerName")

	HTTPGameSocketHandler2(context.Writer, context.Request, gameAPIs.Usecases, id, playerName, hub)
}

func HTTPGameSocketHandler(
	writer http.ResponseWriter,
	request *http.Request,
	usecases *usecases.GameUsecases,
	id int,
	playerName string,
) {
	connection, err := wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}

	game, err := usecases.JoinGame(id, playerName)
	if err != nil {
		if err.Error() != domain.ErrAlreadyInGame {
			err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err))
			if err != nil {
				panic(err)
			}
			connection.Close()
		}
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

		array := strings.Split(message, ": ")
		head := array[0]
		content := strings.Join(array[1:], "/")

		switch head {
		case "leave":
			{
				err = usecases.LeaveGame(id, playerName)
				if err != nil {
					fmt.Println("Could not leave this game: ", err)
					break
				}
				err = SendMessage(connection, "Has left the game")
				if err != nil {
					panic(err)
				}
				connection.Close()
				return
			}
		case "joinTeam":
			{
				err = usecases.JoinTeam(id, playerName, content)
				if err != nil {
					fmt.Println("Could not join this team: ", err)
					break
				}
				game, err := usecases.GetGame(id)
				if err != nil {
					fmt.Println("Could not get updated game: ", err)
					break
				}
				err = sendGame(connection, game)
				if err != nil {
					panic(err)
				}
				break
			}
		default:
			{
				err = SendMessage(connection, "Message not understood by the server")
				if err != nil {
					break
				}
				break
			}
		}

	}
}
