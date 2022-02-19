package api

import (
	"coinche/usecases"
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
	games        map[int]map[*player]bool
	broadcast    chan message
	single       chan private
	register     chan subscription
	unregister   chan subscription
	gameUsecases *usecases.GameUsecases
}

func NewHub(gameUsecases *usecases.GameUsecases) *Hub {
	return &Hub{
		broadcast:    make(chan message),
		single:       make(chan private),
		register:     make(chan subscription),
		unregister:   make(chan subscription),
		games:        make(map[int]map[*player]bool),
		gameUsecases: gameUsecases,
	}
}

func sendToPlayer(connection *websocket.Conn, data []byte) {
	err := send(connection, data)
	if err != nil {
		fmt.Println("Error sending message to player:", err)
	}
}

func deletePlayerAndGameIfNeeded(games map[int]map[*player]bool, players map[*player]bool, player *player, gameID int) {
	close(player.send)
	delete(players, player)
	if len(players) == 0 {
		delete(games, gameID)
	}
}

func register(h *Hub, subscription subscription) {
	players := h.games[subscription.gameID]
	if players == nil {
		players = make(map[*player]bool)
		h.games[subscription.gameID] = players
	}
	h.games[subscription.gameID][subscription.player] = true
}

func unregister(h *Hub, subscription subscription) {
	players := h.games[subscription.gameID]
	if players != nil {
		if _, ok := players[subscription.player]; ok {
			deletePlayerAndGameIfNeeded(h.games, players, subscription.player, subscription.gameID)
		}
	}
}

func broadcast(h *Hub, message message) {
	players := h.games[message.gameID]
	for player := range players {
		select {
		case player.send <- message.data:
			sendToPlayer(player.connection, message.data)
		default:
			deletePlayerAndGameIfNeeded(h.games, players, player, message.gameID)
		}
	}
}

func single(h *Hub, private private) {
	players := h.games[private.gameID]
	player := private.player
	if _, ok := players[player]; !ok {
		data, _ := json.Marshal("Player not in game")
		sendToPlayer(private.player.connection, data)
	}

	select {
	case player.send <- private.data:
		sendToPlayer(player.connection, private.data)
	default:
		deletePlayerAndGameIfNeeded(h.games, players, player, private.gameID)
	}
}

func (h *Hub) run() {
	for {
		select {

		case subscription := <-h.register:
			register(h, subscription)

		case subscription := <-h.unregister:
			unregister(h, subscription)

		case message := <-h.broadcast:
			broadcast(h, message)

		case private := <-h.single:
			single(h, private)
		}
	}
}
