package api

import (
	"coinche/domain"
	"coinche/usecases"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func CreateConnections(test *testing.T, gameUsecases *usecases.GameUsecases, gameID int) (
	*websocket.Conn,
	*websocket.Conn,
	*websocket.Conn,
	*websocket.Conn,
	*httptest.Server,
	*httptest.Server,
	*httptest.Server,
	*httptest.Server,
) {
	var c1 *websocket.Conn
	var s1 *httptest.Server

	var c2 *websocket.Conn
	var s2 *httptest.Server

	var c3 *websocket.Conn
	var s3 *httptest.Server

	var c4 *websocket.Conn
	var s4 *httptest.Server

	hub := NewHub(gameUsecases)
	go hub.run()

	s1, c1 = NewGameWebSocketServer(test, gameID, "P1", &hub)
	s2, c2 = NewGameWebSocketServer(test, gameID, "P2", &hub)
	s3, c3 = NewGameWebSocketServer(test, gameID, "P3", &hub)
	s4, c4 = NewGameWebSocketServer(test, gameID, "P4", &hub)

	_, _ = receive(c1)
	_, _ = receive(c1)
	_, _ = receive(c1)
	_, _ = receive(c1)

	_, _ = receive(c2)
	_, _ = receive(c2)
	_, _ = receive(c2)

	_, _ = receive(c3)
	_, _ = receive(c3)

	_, _ = receive(c4)

	return c1, c2, c3, c4, s1, s2, s3, s4
}

func CloseConnections(
	c1, c2, c3, c4 *websocket.Conn,
	s1, s2, s3, s4 *httptest.Server,
) {
	s1.Close()
	c1.Close()

	s2.Close()
	c2.Close()

	s3.Close()
	c3.Close()

	s4.Close()
	c4.Close()
}

func EmptyMessages(connections []*websocket.Conn, count int) {
	for _, c := range connections {
		for i := 0; i < count; i++ {
			_, err := receive(c)
			if err != nil {
				return
			}
		}
	}
}

func TestSocketTeaming(test *testing.T) {
	assert := assert.New(test)

	gameOne := domain.NewGame("GAME ONE")
	gameOne.Phase = domain.Teaming
	gameOne.ID = 0

	gameTwo := domain.NewGame("GAME TWO")
	gameTwo.Phase = domain.Teaming
	gameTwo.ID = 1

	mockRepository := usecases.NewMockGameRepo(
		map[int]domain.Game{
			0: gameOne,
			1: gameTwo,
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)

	c1, c2, c3, c4, s1, s2, s3, s4 := CreateConnections(test, gameUsecases, 0)

	test.Run("Join a team", func(test *testing.T) {
		err := SendMessage(c1, "joinTeam: AAA", "P1")
		if err != nil {
			test.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond) // prevents concurrent map read and map write

		err = SendMessage(c2, "joinTeam: AAA", "P2")
		if err != nil {
			test.Fatal(err)
		}

		_, _ = receive(c1)

		got := ReceiveGameOrFatal(c1, test)

		EmptyMessages([]*websocket.Conn{c2, c3, c4}, 2)

		assert.Equal("GAME ONE", got.Name)
		assert.Equal("AAA", got.Players["P1"].Team)
		assert.Equal("AAA", got.Players["P2"].Team)
	})

	test.Run("Should fail when joining a team already full", func(test *testing.T) {
		err := SendMessage(c3, "joinTeam: AAA", "P3")
		if err != nil {
			test.Fatal(err)
		}

		reply := ReceiveMessageOrFatal(c3, test)

		assert.Equal("Could not join team: TEAM IS FULL", reply)
	})

	test.Run("Ready to start when two teams ready", func(test *testing.T) {
		err := SendMessage(c3, "joinTeam: BBB", "P3")
		if err != nil {
			test.Fatal(err)
		}

		time.Sleep(50 * time.Millisecond) // prevents concurent map read and map write

		err = SendMessage(c4, "joinTeam: BBB", "P4")
		if err != nil {
			test.Fatal(err)
		}

		_, _ = receive(c1)

		got := ReceiveGameOrFatal(c1, test)

		EmptyMessages([]*websocket.Conn{c2, c3, c4}, 2)

		assert.Equal("GAME ONE", got.Name)
		assert.Equal(32, len(got.Deck))
	})

	test.Run("Can start the game", func(test *testing.T) {
		err := SendMessage(c3, "start", "P3")
		if err != nil {
			test.Fatal(err)
		}

		got := ReceiveGameOrFatal(c3, test)

		EmptyMessages([]*websocket.Conn{c2, c1, c4}, 1)

		assert.Equal(domain.Bidding, got.Phase)
	})

	test.Run("Can place a bid", func(test *testing.T) {
		err := SendMessage(c1, "bid: spade,80", "P1")
		if err != nil {
			test.Fatal(err)
		}

		got := ReceiveGameOrFatal(c1, test)

		EmptyMessages([]*websocket.Conn{c2, c3, c4}, 1)

		assert.Equal("GAME ONE", got.Name)
		assert.Equal(0, got.Bids[80].Coinche)
		assert.Equal(0, got.Bids[80].Pass)
		assert.Equal("P1", got.Bids[80].Player)
		assert.Equal(domain.Spade, got.Bids[80].Color)
	})

	test.Cleanup(func() {
		CloseConnections(c1, c2, c3, c4, s1, s2, s3, s4)
	})
}
