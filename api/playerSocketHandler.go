package api

import (
	"coinche/domain"
	"coinche/usecases"
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	cards = map[string]domain.CardID{
		"7-club":        domain.C7,
		"8-club":        domain.C8,
		"9-club":        domain.C9,
		"10-club":       domain.C10,
		"jack-club":     domain.CJ,
		"queen-club":    domain.CQ,
		"king-club":     domain.CK,
		"as-club":       domain.CA,
		"7-diamond":     domain.D7,
		"8-diamond":     domain.D8,
		"9-diamond":     domain.D9,
		"10-diamond":    domain.D10,
		"jack-diamond":  domain.DJ,
		"queen-diamond": domain.DQ,
		"king-diamond":  domain.DK,
		"as-diamond":    domain.DA,
		"7-heart":       domain.H7,
		"8-heart":       domain.H8,
		"9-heart":       domain.H9,
		"10-heart":      domain.H10,
		"jack-heart":    domain.HJ,
		"queen-heart":   domain.HQ,
		"king-heart":    domain.HK,
		"as-heart":      domain.HA,
		"7-spade":       domain.S7,
		"8-spade":       domain.S8,
		"9-spade":       domain.S9,
		"10-spade":      domain.S10,
		"jack-spade":    domain.SJ,
		"queen-spade":   domain.SQ,
		"king-spade":    domain.SK,
		"as-spade":      domain.SA,
	}
)

func joinGame(connection *websocket.Conn, usecases *usecases.GameUsecases, gameID int, playerName string) domain.Game {
	game, err := usecases.JoinGame(gameID, playerName)
	if err != nil {
		err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err), "S")
		if err != nil {
			fmt.Println("Error sending closing error in « Join game: ", err)
		}
		connection.Close()
		return domain.Game{}
	}

	return game
}

func subscribeAndBroadcast(gameID int, connection *websocket.Conn, game domain.Game, hub *Hub) *player {
	p := &player{hub: hub, connection: connection, send: make(chan []byte, 256)}
	p.hub.register <- subscription{player: p, gameID: gameID}

	broadcastGame(game, p.hub)

	return p
}

type socketHandler struct {
	gameID       int
	playerName   string
	gameUsecases *usecases.GameUsecases
	player       *player
}

func (s *socketHandler) SendErrorMessage(message string, err error) {
	errorMessage := fmt.Sprint(message, err)
	s.player.mu.Lock()
	defer s.player.mu.Unlock()
	err = SendMessage(s.player.connection, errorMessage, "S")
	if err != nil {
		fmt.Println("Error message not sent: ", err)
	}
}

func (s *socketHandler) leave(game domain.Game) {
	err := s.gameUsecases.LeaveGame(s.gameID, s.playerName)
	if err != nil {
		fmt.Println("Could not leave this game: ", err)
		return
	}
	msg := fmt.Sprint(s.playerName, " has left the game")
	broadcastMessage(msg, game.ID, s.player.hub)
	broadcastGame(game, s.player.hub)

	s.player.hub.unregister <- subscription{player: s.player, gameID: s.gameID}

	s.player.connection.Close()
}

func (s *socketHandler) joinTeam(content string) {
	err := s.gameUsecases.JoinTeam(s.gameID, s.playerName, content)
	if err != nil {
		s.SendErrorMessage("Could not join team: ", err)
		return
	}

	game, err := s.gameUsecases.GetGame(s.gameID)
	if err != nil {
		s.SendErrorMessage("Could not get updated game: ", err)
		return
	}

	broadcastGame(game, s.player.hub)
}

func (s socketHandler) startGame(content string) {
	err := s.gameUsecases.StartGame(s.gameID)
	if err != nil {
		s.SendErrorMessage("Could not start game: ", err)
		return
	}

	game, err := s.gameUsecases.GetGame(s.gameID)

	if err != nil {
		s.SendErrorMessage("Could not get updated game: ", err)
		return
	}

	broadcastGame(game, s.player.hub)
}

func (s socketHandler) bid(content string) {
	if content == "pass" {
		err := s.gameUsecases.Pass(s.gameID, s.playerName)
		if err != nil {
			s.SendErrorMessage("Could not pass: ", err)
			return
		}
	} else if content == "coinche" {
		err := s.gameUsecases.Coinche(s.gameID, s.playerName)
		if err != nil {
			s.SendErrorMessage("Could not coinche: ", err)
			return
		}
	} else {
		array := strings.Split(content, ",")
		if len(array) != 2 {
			s.player.mu.Lock()
			defer s.player.mu.Unlock()
			err := SendMessage(s.player.connection, "Invalid bid", "S")
			if err != nil {
				fmt.Println("Error sendind message « Invalid bid » :" + err.Error())
			}
			return
		}

		colorString := array[0]
		valueString := array[1]
		valueInt, err := strconv.Atoi(valueString)
		if err != nil {
			s.SendErrorMessage("Could not parse bid value:", err)
			return
		}

		bidColor := domain.Color(colorString)
		bidValue := domain.BidValue(valueInt)

		err = s.gameUsecases.Bid(s.gameID, s.playerName, bidValue, bidColor)
		if err != nil {
			s.SendErrorMessage("Could not bid: ", err)
			return
		}
	}

	game, err := s.gameUsecases.GetGame(s.gameID)
	if err != nil {
		s.SendErrorMessage("Could not get updated game: ", err)
		return
	}

	broadcastGame(game, s.player.hub)
}

func (s socketHandler) play(content string) {
	card, ok := cards[content]
	if !ok {
		s.player.mu.Lock()
		defer s.player.mu.Unlock()
		err := SendMessage(s.player.connection, "Invalid card", "S")
		if err != nil {
			fmt.Println("Error sending message « Invalid card » : " + err.Error())
		}
		return
	}

	err := s.gameUsecases.PlayCard(s.gameID, s.playerName, card)
	if err != nil {
		s.SendErrorMessage("Could not play: ", err)
		return
	}

	game, err := s.gameUsecases.GetGame(s.gameID)
	if err != nil {
		s.SendErrorMessage("Could not get updated game: ", err)
		return
	}

	broadcastGame(game, s.player.hub)
}

func (s socketHandler) pong() {
	s.player.mu.Lock()
	defer s.player.mu.Unlock()
	err := SendMessageWithoutLog(s.player.connection, "pong")
	if err != nil {
		fmt.Println("Error sending pong message: ", err)
	}
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
			fmt.Println("< Error receiving message: ", err)
			break
		}

		array := strings.Split(message, ": ")
		head := array[0]
		content := strings.Join(array[1:], "/")

		socketHandler := socketHandler{
			gameID:       gameID,
			playerName:   playerName,
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
		case "play":
			{
				socketHandler.play(content)
				break
			}
		case "ping":
			{
				socketHandler.pong()
				break
			}
		default:
			{
				err = SendMessage(connection, "Message not understood by the server", "S")
				if err != nil {
					fmt.Println("Error sending message: « Message not understood by the server » :", err)
				}
				break
			}
		}
	}
}
