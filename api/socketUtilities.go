package api

import (
	"coinche/domain"
	"coinche/utilities"
	testUtilities "coinche/utilities/test"
	"encoding/json"
	"testing"

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

func decodeGame(message []byte) (domain.Game, error) {
	var game domain.Game
	err := json.Unmarshal(message, &game)
	return game, err
}

func decodeMessage(message []byte) (string, error) {
	var reply string
	err := json.Unmarshal(message, &reply)
	return reply, err
}

func ReceiveGameOrFatal(connection *websocket.Conn, test *testing.T) domain.Game {
	message, err := receive(connection)
	testUtilities.FatalIfErr(err, test)

	game, err := decodeGame(message)
	testUtilities.FatalIfErr(err, test)
	return game
}

func ReceiveMessage(connection *websocket.Conn) (string, error) {
	message, err := receive(connection)
	if err != nil {
		return "", err
	}

	reply, err := decodeMessage(message)
	return reply, err
}

func ReceiveMessageOrFatal(connection *websocket.Conn, test *testing.T) string {
	reply, err := ReceiveMessage(connection)
	testUtilities.FatalIfErr(err, test)
	return reply
}
