package api

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"
	"errors"
	"fmt"
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

func Unmarshal(message []byte, destination interface{}) error {
	err := json.Unmarshal(message, &destination)
	if err != nil {
		return errors.New(fmt.Sprint(err, "Message: ", message))
	}
	return nil
}

func decodeGame(message []byte) (domain.Game, error) {
	var game domain.Game
	err := Unmarshal(message, &game)
	if err != nil {
		reply, err := decodeMessage(message)
		if err != nil {
			return game, err
		} else {
			return game, errors.New(fmt.Sprint(err, "Could not decode game: ", reply))
		}
	}
	return game, nil
}

func decodeMessage(message []byte) (string, error) {
	var reply string
	err := Unmarshal(message, &reply)
	return reply, err
}

func ReceiveGameOrFatal(connection *websocket.Conn, test *testing.T) domain.Game {
	message, err := receive(connection)
	if err != nil {
		test.Fatal(err)
	}

	game, err := decodeGame(message)
	if err != nil {
		test.Fatal(err)
	}
	return game
}

func ReceiveMultipleGameOrFatal(c *websocket.Conn, test *testing.T, count int) {
	for i := 0; i < count; i++ {
		_ = ReceiveGameOrFatal(c, test)
	}
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
	if err != nil {
		test.Fatal(err)
	}
	return reply
}
