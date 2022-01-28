package gameapitest

import (
	gameapi "coinche/api/game"
	"encoding/json"
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

func sendMessage(test *testing.T, connection *websocket.Conn, msg string) {
	test.Helper()

	message, err := json.Marshal(msg)
	if err != nil {
		test.Fatal(err)
	}

	if err := connection.WriteMessage(websocket.BinaryMessage, message); err != nil {
		test.Fatalf("%v", err)
	}
}

func receiveWSMessage(test *testing.T, connection *websocket.Conn) string {
	test.Helper()

	_, message, err := connection.ReadMessage()
	if err != nil {
		test.Fatalf("%v", err)
	}

	var reply string
	err = json.Unmarshal(message, &reply)
	if err != nil {
		test.Fatal(err)
	}

	return reply
}

func TestSocketHandler(test *testing.T) {
	assert := assert.New(test)

	test.Run("Can connect and recieve a message", func(test *testing.T) {
		// gin.SetMode(gin.TestMode)

		// w := httptest.NewRecorder()
		// c, _ := gin.CreateTestContext(w)
		// c.Params = []gin.Param{gin.Param{Key: "k", Value: "v"}}

		handler := http.HandlerFunc(gameapi.SocketHandler)

		server, connection := newWSServer(test, handler)
		defer server.Close()
		defer connection.Close()

		sendMessage(test, connection, "hello")

		reply := receiveWSMessage(test, connection)

		assert.Equal("hello", reply)
	})
}
