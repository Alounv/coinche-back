package gameapitest

import (
	gameapi "coinche/api/game"
	"coinche/domain"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func httpToWS(test *testing.T, u string) string {
	test.Helper()

	wsURL, err := url.Parse(u)
	if err != nil {
		test.Fatal(err)
	}

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}

func newWSServer(test *testing.T, handler http.Handler) (*httptest.Server, *websocket.Conn) {
	test.Helper()

	server := httptest.NewServer(handler)
	wsURL := httpToWS(test, server.URL)

	connection, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		test.Fatal(err)
	}

	return server, connection
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

	funcForHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		gameapi.HTTPGameSocketHandler(w, r, &mockUsecases, 1, "player")
	}
	handler := http.HandlerFunc(funcForHandlerFunc)

	server, connection := newWSServer(test, handler)

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
