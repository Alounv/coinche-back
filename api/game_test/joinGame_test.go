package gameapi

import (
	gameapi "coinche/api/game"
	"coinche/domain"
	"coinche/usecases"
	testutils "coinche/utilities/test"
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

	server, connection := testutils.NewGameWebSocketServer(test, gameUsecases, 3, "P1")

	test.Run("Receive error when failing to join", func(test *testing.T) {
		got, _ := gameapi.ReceiveMessage(connection)
		want := "Could not join this game: GAME NOT FOUND"

		assert.Equal(want, got)
	})

	test.Run("Close the connection when failing to join", func(test *testing.T) {
		err := gameapi.SendMessage(connection, "hello")
		if err != nil {
			test.Fatal(err)
		}
		_, err = gameapi.ReceiveMessage(connection)
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

	var s1 *httptest.Server
	var s2 *httptest.Server
	var s3 *httptest.Server
	var s4 *httptest.Server
	var s5 *httptest.Server

	var c1 *websocket.Conn
	var c2 *websocket.Conn
	var c3 *websocket.Conn
	var c4 *websocket.Conn
	var c5 *websocket.Conn

	test.Run("Can connect and receive the game", func(test *testing.T) {
		want := domain.Game(domain.Game{ID: 1, Name: "GAME ONE", Players: map[string]domain.Player{
			"P1": {},
		}})

		s1, c1 = testutils.NewGameWebSocketServer(test, gameUsecases, 1, "P1")

		got, err := gameapi.ReceiveGame(c1)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want, got)
	})

	test.Run("Receive the teaming phase when full", func(test *testing.T) {
		s2, c2 = testutils.NewGameWebSocketServer(test, gameUsecases, 1, "P2")
		s3, c3 = testutils.NewGameWebSocketServer(test, gameUsecases, 1, "P3")
		s4, c4 = testutils.NewGameWebSocketServer(test, gameUsecases, 1, "P4")

		got, err := gameapi.ReceiveGame(c4)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal("GAME ONE", got.Name)
		assert.Equal(map[string]domain.Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}}, got.Players)
		assert.Equal(domain.Teaming, got.Phase)
	})

	test.Run("Try to join when already in game", func(test *testing.T) {
		s5, c5 = testutils.NewGameWebSocketServer(test, gameUsecases, 1, "P4")

		got, err := gameapi.ReceiveGame(c5)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal("GAME ONE", got.Name)
		assert.Equal(map[string]domain.Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}}, got.Players)
		assert.Equal(domain.Teaming, got.Phase)
	})

	test.Run("Can send a message", func(test *testing.T) {
		err := gameapi.SendMessage(c1, "hello")
		if err != nil {
			test.Fatal(err)
		}
		reply, err := gameapi.ReceiveMessage(c1)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal("hello", reply)
	})

	test.Run("Can leave the game", func(test *testing.T) {
		err := gameapi.SendMessage(c1, "leave")
		if err != nil {
			test.Fatal(err)
		}
		reply, err := gameapi.ReceiveMessage(c1)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal("Has left the game", reply)
	})

	test.Run("Can close the connection", func(test *testing.T) {
		c1.Close()
		err := gameapi.SendMessage(c1, "hello")

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
