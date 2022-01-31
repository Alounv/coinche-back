package testutils

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

func NewWSServer(test *testing.T, handler http.Handler) (*httptest.Server, *websocket.Conn) {
	test.Helper()

	server := httptest.NewServer(handler)
	wsURL := httpToWS(test, server.URL)

	connection, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		test.Fatal(err)
	}

	return server, connection
}
