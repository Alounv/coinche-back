package gameapitest

import (
	gameapi "coinche/api/game"
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
	handler := http.HandlerFunc(gameapi.HTTPGameSocketHandler)

	server, connection := newWSServer(test, handler)

	test.Run("Can connect and receive a message", func(test *testing.T) {
		reply1, _ := gameapi.ReceiveMessage(connection)

		assert.Equal("connection established", reply1)
	})

	test.Run("Can send a message", func(test *testing.T) {
		err := gameapi.SendMessage(connection, "hello")
		if err != nil {
			test.Fatal(err)
		}
		reply, _ := gameapi.ReceiveMessage(connection)

		assert.Equal("hello", reply)
	})

	/*test.Run("Can close the connection", func(test *testing.T) {
		connection.Close()
		err := gameapi.SendMessage(connection, "hello")

		assert.EqualError(err, "write tcp 127.0.0.1:52624->127.0.0.1:52623: use of closed network connection")
	})*/

	test.Cleanup(func() {
		server.Close()
		connection.Close()
	})
}
