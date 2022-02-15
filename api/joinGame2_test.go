package api

import (
	"coinche/domain"
	"coinche/usecases"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSocketHandler2(test *testing.T) {
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

	/*var s2 *httptest.Server
	var s3 *httptest.Server
	var s4 *httptest.Server
	var s5 *httptest.Server

	var c2 *websocket.Conn
	var c3 *websocket.Conn
	var c4 *websocket.Conn
	var c5 *websocket.Conn*/

	hub := newHub()
	go hub.run()

	test.Run("Can connect and receive the game", func(test *testing.T) {
		want := domain.Game(domain.Game{ID: 1, Name: "GAME ONE", Players: map[string]domain.Player{
			"P1": {},
		}})

		s1, c1 = NewGameWebSocketServer2(test, gameUsecases, 1, "P1", hub)

		message, err := receive(c1)
		if err != nil {
			test.Fatal(err)
		}

		got, err := DecodeGame(message)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want, got)
	})

	test.Cleanup(func() {
		s1.Close()
		c1.Close()
		/*	s2.Close()
			c2.Close()
			s3.Close()
			c3.Close()
			s4.Close()
			c4.Close()
			s5.Close()
			c5.Close()*/
	})
}
