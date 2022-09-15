package api

import (
	"coinche/domain"
	"coinche/usecases"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func CreateConnections(test *testing.T, gameUsecases *usecases.GameUsecases) (
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

	s1, c1 = NewGameWebSocketServer(test, 1, "P1", &hub)
	s2, c2 = NewGameWebSocketServer(test, 1, "P2", &hub)
	s3, c3 = NewGameWebSocketServer(test, 1, "P3", &hub)
	s4, c4 = NewGameWebSocketServer(test, 1, "P4", &hub)

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
	mockRepository := usecases.NewMockGameRepo(
		map[int]domain.Game{
			1: {ID: 1, Name: "GAME ONE", Phase: domain.Preparation, Players: map[string]domain.Player{}, Deck: domain.NewDeck()},
			2: {ID: 2, Name: "GAME TWO", Phase: domain.Preparation, Players: map[string]domain.Player{}, Deck: domain.NewDeck()},
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)

	c1, c2, c3, c4, s1, s2, s3, s4 := CreateConnections(test, gameUsecases)

	test.Run("Join a team", func(test *testing.T) {
		err := SendMessage(c1, "joinTeam: AAA")
		if err != nil {
			test.Fatal(err)
		}

		err = SendMessage(c2, "joinTeam: AAA")
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
		err := SendMessage(c3, "joinTeam: AAA")
		if err != nil {
			test.Fatal(err)
		}

		reply := ReceiveMessageOrFatal(c3, test)

		assert.Equal("Could not join this team: TEAM IS FULL", reply)
	})

	test.Run("Ready to start when two teams ready", func(test *testing.T) {
		err := SendMessage(c3, "joinTeam: BBB")
		if err != nil {
			test.Fatal(err)
		}

		err = SendMessage(c4, "joinTeam: BBB")
		if err != nil {
			test.Fatal(err)
		}

		_, _ = receive(c1)

		got := ReceiveGameOrFatal(c1, test)

		EmptyMessages([]*websocket.Conn{c2, c3, c4}, 2)

		assert.Equal("GAME ONE", got.Name)
		assert.NoError(got.CanStartBidding())
		assert.Equal(32, len(got.Deck))
	})

	test.Run("Can start the game", func(test *testing.T) {
		err := SendMessage(c3, "start")
		if err != nil {
			test.Fatal(err)
		}

		got := ReceiveGameOrFatal(c3, test)

		assert.Equal(domain.Bidding, got.Phase)
	})

	// TODO: TEST PLACE SOME BIDS

	// TODO: TEST PLAY CARDS

	// TODO: TEST COUNTING

	// TODO: TEST RESTART

	// TODO: TEST SCORES ON MULTIPLE

	test.Cleanup(func() {
		CloseConnections(c1, c2, c3, c4, s1, s2, s3, s4)
	})
}
