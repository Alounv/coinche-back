package api

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"

	"github.com/gorilla/websocket"
)

func send(connection *websocket.Conn, message []byte) error {
	err := connection.WriteMessage(websocket.BinaryMessage, message)
	return err
}

func receive(connection *websocket.Conn) ([]byte, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, err
}

func broadcastGameOrPanic(game domain.Game, hub *Hub) {
	data, err := json.Marshal(game)
	utilities.PanicIfErr(err)

	m := message{data: data, gameID: game.ID}

	hub.broadcast <- m
}

func SendMessage(connection *websocket.Conn, msg string) error {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = send(connection, message)
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
