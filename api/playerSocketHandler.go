package api

import (
	"coinche/domain"
	"coinche/usecases"
	"coinche/utilities"
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

func joinGame(connection *websocket.Conn, usecases *usecases.GameUsecases, gameID int, playerName string) domain.Game {
	game, err := usecases.JoinGame(gameID, playerName)
	if err != nil {
		err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err))
		utilities.PanicIfErr(err)
		connection.Close()
		return domain.Game{}
	}

	return game
}

func subscribeAndBroadcast(gameID int, connection *websocket.Conn, game domain.Game, hub *Hub) *player {
	p := &player{hub: hub, connection: connection, send: make(chan []byte, 256)}
	p.hub.register <- subscription{player: p, gameID: gameID}

	broadcastGameOrPanic(game, p.hub)

	return p
}

type socketHandler struct {
	gameID       int
	playerName   string
	connection   *websocket.Conn
	gameUsecases *usecases.GameUsecases
	player       *player
}

func (s *socketHandler) SendErrorMessageOrPanic(message string, err error) {
	errorMessage := fmt.Sprint(message, err)
	err = SendMessage(s.connection, errorMessage)
	utilities.PanicIfErr(err)
}

func (s *socketHandler) leave(game domain.Game) {
	err := s.gameUsecases.LeaveGame(s.gameID, s.playerName)
	if err != nil {
		fmt.Println("Could not leave this game: ", err)
		return
	}
	err = SendMessage(s.connection, "Has left the game")
	utilities.PanicIfErr(err)
	broadcastGameOrPanic(game, s.player.hub)

	s.player.hub.unregister <- subscription{player: s.player, gameID: s.gameID}

	s.connection.Close()
}

func (s *socketHandler) joinTeam(content string) {
	err := s.gameUsecases.JoinTeam(s.gameID, s.playerName, content)
	if err != nil {
		s.SendErrorMessageOrPanic("Could not join team: ", err)
		return
	}

	game, err := s.gameUsecases.GetGame(s.gameID)
	if err != nil {
		s.SendErrorMessageOrPanic("Could not get updated game: ", err)
		return
	}

	broadcastGameOrPanic(game, s.player.hub)
}

func (s socketHandler) startGame(content string) {
	err := s.gameUsecases.StartGame(s.gameID)
	if err != nil {
		s.SendErrorMessageOrPanic("Could not start game: ", err)
		return
	}

	game, err := s.gameUsecases.GetGame(s.gameID)

	if err != nil {
		s.SendErrorMessageOrPanic("Could not get updated game: ", err)
		return
	}

	broadcastGameOrPanic(game, s.player.hub)
}

func (s socketHandler) bid(content string) {
	if content == "pass" {
		err := s.gameUsecases.Pass(s.gameID, s.playerName)
		if err != nil {
			s.SendErrorMessageOrPanic("Could not pass: ", err)
			return
		}
	} else if content == "coinche" {
		err := s.gameUsecases.Coinche(s.gameID, s.playerName)
		if err != nil {
			s.SendErrorMessageOrPanic("Could not coinche: ", err)
			return
		}
	} else {

		array := strings.Split(content, ",")
		if len(array) != 2 {
			err := SendMessage(s.connection, "Invalid bid")
			utilities.PanicIfErr(err)
			return
		}

		colorString := array[0]
		valueString := array[1]
		valueInt, err := strconv.Atoi(valueString)
		if err != nil {
			s.SendErrorMessageOrPanic("Could not parse bid value:", err)
			return
		}

		bidColor := domain.Color(colorString)
		bidValue := domain.BidValue(valueInt)

		err = s.gameUsecases.Bid(s.gameID, s.playerName, bidValue, bidColor)
		if err != nil {
			s.SendErrorMessageOrPanic("Could not bid: ", err)
			return
		}
	}

	game, err := s.gameUsecases.GetGame(s.gameID)
	if err != nil {
		s.SendErrorMessageOrPanic("Could not get updated game: ", err)
		return
	}

	broadcastGameOrPanic(game, s.player.hub)
}

func PlayerSocketHandler(
	connection *websocket.Conn,
	gameID int,
	playerName string,
	hub *Hub,
) {
	game := joinGame(connection, hub.gameUsecases, gameID, playerName)
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
			gameID:       gameID,
			playerName:   playerName,
			connection:   connection,
			gameUsecases: hub.gameUsecases,
			player:       player,
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
		case "bid":
			{
				socketHandler.bid(content)
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
