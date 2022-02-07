package gameapitest

import (
	gameapi "coinche/api/game"
	"coinche/domain"
	testutils "coinche/utilities/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailingSocketHandler(test *testing.T) {
	assert := assert.New(test)
	mockUsecases := MockGameUsecases{
		map[int]domain.Game{},
		nil,
	}

	server, connection := testutils.NewGameWebSocketServer(test, &mockUsecases, 1, "player")

	test.Run("Receive error when failing to join", func(test *testing.T) {
		got, _ := gameapi.ReceiveMessage(connection)
		want := "Could not join this game: TEST JOIN FAIL"

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
	mockUsecases := MockGameUsecases{
		map[int]domain.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO"},
		},
		nil,
	}

	server, connection := testutils.NewGameWebSocketServer(test, &mockUsecases, 1, "player")

	test.Run("Can connect and receive the game", func(test *testing.T) {
		got, _ := gameapi.ReceiveGame(connection)
		want := domain.Game(domain.Game{ID: 1, Name: "GAME ONE", Players: []string{"player"}})

		assert.Equal(want, got)
	})

	test.Run("Can send a message", func(test *testing.T) {
		err := gameapi.SendMessage(connection, "hello")
		if err != nil {
			test.Fatal(err)
		}
		reply, _ := gameapi.ReceiveMessage(connection)

		assert.Equal("hello", reply)
	})

	test.Run("Can leave the game", func(test *testing.T) {
		err := gameapi.SendMessage(connection, "leave")
		if err != nil {
			test.Fatal(err)
		}
		reply, _ := gameapi.ReceiveMessage(connection)

		assert.Equal("Has left the game", reply)
	})

	test.Run("Can close the connection", func(test *testing.T) {
		connection.Close()
		err := gameapi.SendMessage(connection, "hello")

		assert.NotNil(err)
	})

	test.Cleanup(func() {
		server.Close()
		connection.Close()
	})
}
