package api

import (
	"coinche/domain"
	"coinche/usecases"
	"coinche/utilities"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestFailingSocketHandler(test *testing.T) {
	assert := assert.New(test)
	mockRepository := usecases.NewMockGameRepo(
		map[int]*domain.Game{
			1: {ID: 1, Name: "GAME ONE", Phase: domain.Preparation, Players: map[string]domain.Player{}},
			2: {ID: 2, Name: "GAME TWO", Phase: domain.Preparation, Players: map[string]domain.Player{}},
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)

	hub := NewHub()
	go hub.run()
	server, connection := NewGameWebSocketServer(test, gameUsecases, 3, "P1", hub)

	test.Run("Receive error when failing to join", func(test *testing.T) {
		got, _ := ReceiveMessage(connection)
		want := "Could not join this game: GAME NOT FOUND"

		assert.Equal(want, got)
	})

	test.Run("Close the connection when failing to join", func(test *testing.T) {
		err := SendMessage(connection, "hello")
		utilities.FatalIfErr(err, test)
		_, err = ReceiveMessage(connection)
		assert.NotNil(err)
	})

	test.Cleanup(func() {
		server.Close()
		connection.Close()
	})
}

func TestSocketHandler(test *testing.T) {
	assert := assert.New(test)

	mockRepository := usecases.NewMockGameRepo(
		map[int]*domain.Game{
			1: {ID: 1, Name: "GAME ONE", Phase: domain.Preparation, Players: map[string]domain.Player{}},
			2: {ID: 2, Name: "GAME TWO", Phase: domain.Preparation, Players: map[string]domain.Player{}},
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)
	var c1 *websocket.Conn
	var s1 *httptest.Server

	var c2 *websocket.Conn
	var s2 *httptest.Server

	var c3 *websocket.Conn
	var s3 *httptest.Server

	var c4 *websocket.Conn
	var s4 *httptest.Server

	var s5 *httptest.Server
	var c5 *websocket.Conn

	hub := NewHub()
	go hub.run()

	test.Run("Can connect and receive the game", func(test *testing.T) {
		want := domain.Game(domain.Game{ID: 1, Name: "GAME ONE", Players: map[string]domain.Player{
			"P1": {},
		}})

		s1, c1 = NewGameWebSocketServer(test, gameUsecases, 1, "P1", hub)

		message, err := receive(c1)
		utilities.FatalIfErr(err, test)

		got, err := DecodeGame(message)
		utilities.FatalIfErr(err, test)

		assert.Equal(want, got)
	})

	test.Run("Receive the teaming phase when full", func(test *testing.T) {
		s2, c2 = NewGameWebSocketServer(test, gameUsecases, 1, "P2", hub)
		s3, c3 = NewGameWebSocketServer(test, gameUsecases, 1, "P3", hub)
		s4, c4 = NewGameWebSocketServer(test, gameUsecases, 1, "P4", hub)

		_, _ = receive(c1)
		_, _ = receive(c1)
		message, err := receive(c1)
		utilities.FatalIfErr(err, test)

		_, _ = receive(c2)
		_, _ = receive(c2)
		_, _ = receive(c2)

		_, _ = receive(c3)
		_, _ = receive(c3)

		_, _ = receive(c4)

		got, err := DecodeGame(message)
		utilities.FatalIfErr(err, test)

		assert.Equal("GAME ONE", got.Name)
		assert.Equal(map[string]domain.Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}}, got.Players)
		assert.Equal(domain.Teaming, got.Phase)
	})

	test.Run("Try to join when already in game", func(test *testing.T) {
		s5, c5 = NewGameWebSocketServer(test, gameUsecases, 1, "P4", hub)

		got, err := ReceiveMessage(c5)
		utilities.FatalIfErr(err, test)

		assert.Equal("Could not join this game: ALREADY IN GAME", got)
	})

	test.Run("Try to join a full game", func(test *testing.T) {
		s5, c5 = NewGameWebSocketServer(test, gameUsecases, 1, "P5", hub)

		message, err := receive(c5)
		utilities.FatalIfErr(err, test)

		got, err := DecodeMessage(message)
		utilities.FatalIfErr(err, test)

		assert.Equal("Could not join this game: GAME IS FULL", got)
	})

	test.Run("Can send a message", func(test *testing.T) {
		err := SendMessage(c1, "hello")
		utilities.FatalIfErr(err, test)
		reply, err := ReceiveMessage(c1)
		utilities.FatalIfErr(err, test)

		assert.Equal("Message not understood by the server", reply)
	})

	test.Run("Can leave the game", func(test *testing.T) {
		err := SendMessage(c1, "leave")
		utilities.FatalIfErr(err, test)
		reply, err := ReceiveMessage(c1)
		utilities.FatalIfErr(err, test)

		assert.Equal("Has left the game", reply)
	})

	test.Run("Can close the connection", func(test *testing.T) {
		c1.Close()
		err := SendMessage(c1, "hello")

		assert.NotNil(err)
	})

	test.Cleanup(func() {
		s1.Close()
		c1.Close()

		s2.Close()
		c2.Close()

		s3.Close()
		c3.Close()

		s4.Close()
		c4.Close()

		s5.Close()
		c5.Close()
	})
}
