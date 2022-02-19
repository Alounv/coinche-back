package api

import (
	"coinche/domain"
	"coinche/usecases"
	"coinche/utilities"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func joinGame(connection *websocket.Conn, usecases *usecases.GameUsecases, id int, playerName string) domain.Game {
	game, err := usecases.JoinGame(id, playerName)
	if err != nil {
		fmt.Println("Error joining game:", err)
		err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err))
		utilities.PanicIfErr(err)
		connection.Close()
		return domain.Game{}
	}

	return game
}

func subscribeAndBroadcast(id int, connection *websocket.Conn, game domain.Game, hub *Hub) *player {
	p := &player{hub: hub, connection: connection, send: make(chan []byte, 256)}
	p.hub.register <- subscription{player: p, gameID: id}

	broadcastGameOrPanic(game, p.hub)

	return p
}

type socketHandler struct {
	gameID     int
	playerName string
	connection *websocket.Conn
	usecases   *usecases.GameUsecases
	player     *player
}

func (s *socketHandler) leave(game domain.Game) {
	err := s.usecases.LeaveGame(s.gameID, s.playerName)
	if err != nil {
		fmt.Println("Could not leave this game: ", err)
		return
	}
	err = SendMessage(s.connection, "Has left the game")
	utilities.PanicIfErr(err)
	broadcastGameOrPanic(game, s.player.hub)

	s.connection.Close()
}

func (s *socketHandler) joinTeam(content string) {
	err := s.usecases.JoinTeam(s.gameID, s.playerName, content)
	if err != nil {
		errorMessage := fmt.Sprint("Could not join this team: ", err)
		err = SendMessage(s.connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	game, err := s.usecases.GetGame(s.gameID)
	if err != nil {
		errorMessage := fmt.Sprint("Could not get updated game: ", err)
		err := SendMessage(s.connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	broadcastGameOrPanic(game, s.player.hub)
}

func (s socketHandler) startGame(content string) {
	err := s.usecases.StartGame(s.gameID)
	if err != nil {
		errorMessage := fmt.Sprint("Could not start the game: ", err)
		err = SendMessage(s.connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	game, err := s.usecases.GetGame(s.gameID)
	if err != nil {
		errorMessage := fmt.Sprint("Could not get updated game: ", err)
		err := SendMessage(s.connection, errorMessage)
		utilities.PanicIfErr(err)
		return
	}
	broadcastGameOrPanic(game, s.player.hub)
}

func PlayerSocketHandler(
	connection *websocket.Conn,
	usecases *usecases.GameUsecases,
	gameID int,
	playerName string,
	hub *Hub,
) {
	game := joinGame(connection, usecases, gameID, playerName)
	player := subscribeAndBroadcast(gameID, connection, game, hub)

	for {
		message, err := ReceiveMessage(connection)
		if err != nil {
			break
		}

		array := strings.Split(message, ": ")
		head := array[0]
		content := strings.Join(array[1:], "/")

		socketHandler := socketHandler{
			gameID:     gameID,
			playerName: playerName,
			connection: connection,
			usecases:   usecases,
			player:     player,
		}

		switch head {
		case "leave":
			{
				socketHandler.leave(game)
				break
			}
		case "joinTeam":
			{
				socketHandler.joinTeam(content)
				break
			}
		case "start":
			{
				socketHandler.startGame(content)
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
