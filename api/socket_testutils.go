package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
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

func newConnection(test *testing.T, serverURL string) *websocket.Conn {
	wsURL := httpToWS(test, serverURL)

	connection, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		test.Fatal(err)
	}
	return connection
}

func newServer(test *testing.T, handler http.Handler) (*httptest.Server, *websocket.Conn) {
	test.Helper()

	server := httptest.NewServer(handler)

	connection := newConnection(test, server.URL)

	return server, connection
}

func NewGameWebSocketServer(
	test *testing.T,
	ID int,
	playerName string,
	hub *Hub,
) (*httptest.Server, *websocket.Conn) {
	funcForHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		connection, err := wsupgrader.Upgrade(w, r, nil)
		if err != nil {
			test.Fatal(err)
		}
		PlayerSocketHandler(connection, ID, playerName, hub)
	}
	socketHandler := http.HandlerFunc(funcForHandlerFunc)

	return newServer(test, socketHandler)
}
