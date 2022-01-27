package ports

import (
	"coinche/app"
	testUtils "coinche/utilities/test"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinGame(test *testing.T) {
	assert := assert.New(test)
	mockGameService := MockGameService{
		map[int]app.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO", Full: true},
		},
		nil,
	}
	router := SetupRouter(&mockGameService)

	test.Run("join game 1", func(test *testing.T) {
		request := NewJoinGameRequest(1, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusAccepted, response.Code)
	})

	test.Run("should fail when joining non existing game", func(test *testing.T) {
		request := NewJoinGameRequest(60, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusNotFound, response.Code)
	})

	test.Run("should fail when game is full", func(test *testing.T) {
		request := NewJoinGameRequest(2, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusForbidden, response.Code)
	})
}

func NewJoinGameRequest(id int, playerName string) *http.Request {
	route := fmt.Sprintf("/games/%d/join?playerName=%s", id, url.QueryEscape(playerName))
	return testUtils.GetNewRequest(route, http.MethodPost)
}
