package api

import (
	"coinche/domain"
	"coinche/usecases"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type player struct {
	hub        *hub
	connection *websocket.Conn
	send       chan []byte
}

type message struct {
	data   []byte
	gameID int
}

type private struct {
	player *player
	data   []byte
	gameID int
}

type subscription struct {
	player *player
	gameID int
}

type hub struct {
	games      map[int]map[*player]bool
	broadcast  chan message
	single     chan private
	register   chan subscription
	unregister chan subscription
}

func newHub() *hub {
	return &hub{
		broadcast:  make(chan message),
		single:     make(chan private),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		games:      make(map[int]map[*player]bool),
	}
}

func (h *hub) run() {
	for {
		select {

		case subscription := <-h.register:
			players := h.games[subscription.gameID]
			if players == nil {
				players = make(map[*player]bool)
				h.games[subscription.gameID] = players
			}
			h.games[subscription.gameID][subscription.player] = true

		case subscription := <-h.unregister:
			players := h.games[subscription.gameID]
			if players != nil {
				if _, ok := players[subscription.player]; ok {
					delete(players, subscription.player)
					close(subscription.player.send)
					if len(players) == 0 {
						delete(h.games, subscription.gameID)
					}
				}
			}

		case message := <-h.broadcast:
			players := h.games[message.gameID]
			for player := range players {
				select {
				case player.send <- message.data:
					err := send(player.connection, message.data)
					if err != nil {
						fmt.Println("Error sending message to player:", err)
					}
				default:
					close(player.send)
					delete(players, player)
					if len(players) == 0 {
						delete(h.games, message.gameID)
					}
				}
			}

		case private := <-h.single:

			players := h.games[private.gameID]
			player := private.player
			if _, ok := players[player]; !ok {
				message, _ := json.Marshal("Player not in game")
				err := send(private.player.connection, message)
				if err != nil {
					fmt.Println("Error sending message to player:", err)
				}
			}

			select {
			case player.send <- private.data:
				err := send(player.connection, private.data)
				if err != nil {
					fmt.Println("Error sending message to player:", err)
				}
			default:
				close(player.send)
				delete(players, player)
				if len(players) == 0 {
					delete(h.games, private.gameID)
				}
			}
		}
	}
}

func HTTPGameSocketHandler2(
	writer http.ResponseWriter,
	request *http.Request,
	usecases *usecases.GameUsecases,
	id int,
	playerName string,
	hub *hub,
) {
	connection, err := wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}

	game, err := usecases.JoinGame(id, playerName)
	if err != nil {
		fmt.Println("Error joining game:", err)
		err := SendMessage(connection, fmt.Sprint("Could not join this game: ", err))
		if err != nil {
			panic(err)
		}
		connection.Close()
		return
	}

	p := &player{hub: hub, connection: connection, send: make(chan []byte, 256)}
	p.hub.register <- subscription{player: p, gameID: id}

	err = broadcastGame(game, p.hub)
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
				err = broadcastGame(game, p.hub)
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
				err = broadcastGame(game, p.hub)
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

func broadcastGame(game domain.Game, hub *hub) error {
	data, err := json.Marshal(game)
	if err != nil {
		return err
	}

	m := message{data: data, gameID: game.ID}

	hub.broadcast <- m
	return nil
}

/*func sendPrivateMessage(message string, toPlayer *player, gameID int, hub *hub) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m := private{data: data, gameID: gameID, player: toPlayer}

	hub.single <- m
	return nil
}*/

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

func DecodeGame(message []byte) (domain.Game, error) {
	var game domain.Game
	err := json.Unmarshal(message, &game)
	return game, err
}

func DecodeMessage(message []byte) (string, error) {
	var reply string
	err := json.Unmarshal(message, &reply)
	return reply, err
}

func ReceiveGame(connection *websocket.Conn) (domain.Game, error) {
	message, err := receive(connection)
	if err != nil {
		return domain.Game{}, err
	}

	game, err := DecodeGame(message)
	return game, err
}

func ReceiveMessage(connection *websocket.Conn) (string, error) {
	message, err := receive(connection)
	if err != nil {
		return "", err
	}

	reply, err := DecodeMessage(message)
	return reply, err
}

func receive(connection *websocket.Conn) ([]byte, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, err
}
