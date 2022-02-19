package api

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type player struct {
	hub        *Hub
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

type Hub struct {
	games      map[int]map[*player]bool
	broadcast  chan message
	single     chan private
	register   chan subscription
	unregister chan subscription
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan message),
		single:     make(chan private),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		games:      make(map[int]map[*player]bool),
	}
}

func (h *Hub) run() {
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
