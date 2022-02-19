package api

import (
	"coinche/usecases"
	"coinche/utilities"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func httpToWS(test *testing.T, u string) string {
	test.Helper()

	wsURL, err := url.Parse(u)
	utilities.FatalIfErr(err, test)

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}

func newConnection(test *testing.T, wsURL string) *websocket.Conn {
	connection, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	utilities.FatalIfErr(err, test)
	return connection
}

func newServer(test *testing.T, handler http.Handler) (*httptest.Server, *websocket.Conn) {
	test.Helper()

	server := httptest.NewServer(handler)
	wsURL := httpToWS(test, server.URL)

	connection := newConnection(test, wsURL)

	return server, connection
}

func NewGameWebSocketServer(
	test *testing.T,
	gameUsecases *usecases.GameUsecases,
	ID int,
	playerName string,
	hub *Hub,
) (*httptest.Server, *websocket.Conn) {
	funcForHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		connection, err := wsupgrader.Upgrade(w, r, nil)
		utilities.FatalIfErr(err, test)
		PlayerSocketHandler(connection, gameUsecases, ID, playerName, hub)
	}
	socketHandler := http.HandlerFunc(funcForHandlerFunc)

	return newServer(test, socketHandler)
}
