package api

import (
	"coinche/usecases"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type player struct {
	hub        *Hub
	connection *websocket.Conn
	send       chan []byte
	mu         sync.Mutex
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

func NewHub(gameUsecases *usecases.GameUsecases) Hub {
	return Hub{
		broadcast:    make(chan message),
		single:       make(chan private),
		register:     make(chan subscription),
		unregister:   make(chan subscription),
		games:        make(map[int]map[*player]bool),
		gameUsecases: gameUsecases,
	}
}

func sendToPlayerOrUnregister(h *Hub, player *player, data []byte, gameID int) {
	player.mu.Lock()
	defer player.mu.Unlock()
	err := send(player.connection, data)
	if err != nil {
		fmt.Println("Error sending message to player: (", err, "). Closing connection")
		deletePlayerAndGameIfNeeded(h.games, h.games[gameID], player, gameID)
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
			sendToPlayerOrUnregister(h, player, message.data, message.gameID)
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
		sendToPlayerOrUnregister(h, private.player, data, private.gameID)
	}

	select {
	case player.send <- private.data:
		sendToPlayerOrUnregister(h, player, private.data, private.gameID)
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
