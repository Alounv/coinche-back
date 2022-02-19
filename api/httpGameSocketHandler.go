package api

import (
	"coinche/domain"
	"coinche/usecases"
	"coinche/utilities"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func setup(writer http.ResponseWriter,
	request *http.Request,
	usecases *usecases.GameUsecases,
	id int,
	playerName string,
	hub *Hub) (*websocket.Conn, *player, domain.Game) {
	connection, err := wsupgrader.Upgrade(writer, request, nil)
	utilities.PanicIfErr(err)

	game, err := usecases.JoinGame(id, playerName)
	if err != nil {
		fmt.Println("Error joining game:", err)
		err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err))
		utilities.PanicIfErr(err)
		connection.Close()
		return nil, nil, domain.Game{}
	}

	p := &player{hub: hub, connection: connection, send: make(chan []byte, 256)}
	p.hub.register <- subscription{player: p, gameID: id}

	broadcastGameOrPanic(game, p.hub)
	utilities.PanicIfErr(err)

	return connection, p, game
}

func leave(id int, playerName string, connection *websocket.Conn, usecases *usecases.GameUsecases, game domain.Game, player *player) {
	err := usecases.LeaveGame(id, playerName)
	if err != nil {
		fmt.Println("Could not leave this game: ", err)
		return
	}
	err = SendMessage(connection, "Has left the game")
	utilities.PanicIfErr(err)
	broadcastGameOrPanic(game, player.hub)

	connection.Close()
}

func joinTeam(id int, playerName string, connection *websocket.Conn, usecases *usecases.GameUsecases, game domain.Game, player *player, content string) {
	err := usecases.JoinTeam(id, playerName, content)
	if err != nil {
		errorMessage := fmt.Sprint("Could not join this team: ", err)
		err = SendMessage(connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	game, err = usecases.GetGame(id)
	if err != nil {
		errorMessage := fmt.Sprint("Could not get updated game: ", err)
		err := SendMessage(connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	broadcastGameOrPanic(game, player.hub)
}

func startGame(id int, playerName string, connection *websocket.Conn, usecases *usecases.GameUsecases, game domain.Game, player *player, content string) {
	err := usecases.StartGame(id)
	if err != nil {
		errorMessage := fmt.Sprint("Could not start the game: ", err)
		err = SendMessage(connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	game, err = usecases.GetGame(id)
	if err != nil {
		errorMessage := fmt.Sprint("Could not get updated game: ", err)
		err := SendMessage(connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	broadcastGameOrPanic(game, player.hub)
}

func HTTPGameSocketHandler(
	writer http.ResponseWriter,
	request *http.Request,
	usecases *usecases.GameUsecases,
	id int,
	playerName string,
	hub *Hub,
) {
	connection, player, game := setup(writer, request, usecases, id, playerName, hub)
	if connection == nil {
		return
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
				leave(id, playerName, connection, usecases, game, player)
				break
			}
		case "joinTeam":
			{
				joinTeam(id, playerName, connection, usecases, game, player, content)
				break
			}
		case "start":
			{
				startGame(id, playerName, connection, usecases, game, player, content)
				break
			}
		default:
			{
				err = SendMessage(connection, "Message not understood by the server")
				utilities.PanicIfErr(err)
				break
			}
		}
	}
}
