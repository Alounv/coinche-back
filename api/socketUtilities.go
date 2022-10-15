package api

import (
	"coinche/domain"
	testUtilities "coinche/utilities/test"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"testing"
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

func broadcastGame(game domain.Game, hub *Hub) {
	fmt.Println("S >>> broadcasting game:", game.ID)
	data, err := json.Marshal(game)
	if err != nil {
		fmt.Println("Error marshal during broadcasting game: " + err.Error())
		return
	}

	m := message{data: data, gameID: game.ID}

	hub.broadcast <- m
}

func broadcastMessage(msg string, gameID int, hub *Hub) {
	fmt.Println("S >>> broadcasting message:", msg)
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshal during broadcasting message: " + err.Error())
		return
	}

	m := message{data: data, gameID: gameID}

	hub.broadcast <- m
}

func SendMessageOrFatal(connection *websocket.Conn, msg string, origin string, test *testing.T) {
	err := SendMessage(connection, msg, origin)
	if err != nil {
		test.Fatal(err)
	}
}

func SendMessage(connection *websocket.Conn, msg string, origin string) error {
	fmt.Println(origin, "> sending message:", msg)
	err := SendMessageWithoutLog(connection, msg)
	return err
}

func SendMessageWithoutLog(connection *websocket.Conn, msg string) error {
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
		}

		return game, errors.New(fmt.Sprint(err, "Could not decode game: ", reply))
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

func EmtpyGamesForOtherPlayersOrFatal(names []string, currentPlayerName string, gameID int, test *testing.T, connections []*websocket.Conn) {
	for _, name := range names {
		p := testUtilities.GetPlayerIndexFromNameOrFatal(name, test)
		if name == currentPlayerName {
			continue
		}
		ReceiveMultipleGameOrFatal(connections[p], test, gameID)
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

func ReceiveMessageOrGame(connection *websocket.Conn) (string, domain.Game, error) {
	message, err := receive(connection)
	if err != nil {
		return "", domain.Game{}, err
	}

	reply, err := decodeMessage(message)
	if err == nil {
		return reply, domain.Game{}, nil
	}

	game, err := decodeGame(message)
	if err == nil {
		return "", game, nil
	}

	return "", domain.Game{}, err
}

func ReceiveMessageOrFatal(connection *websocket.Conn, test *testing.T) string {
	reply, err := ReceiveMessage(connection)
	if err != nil {
		test.Fatal(err)
	}
	return reply
}

func ReceiveMessageOrGameOrFatal(connection *websocket.Conn, test *testing.T) (string, domain.Game) {
	reply, game, err := ReceiveMessageOrGame(connection)
	if err != nil {
		test.Fatal(err)
	}
	return reply, game
}
