package gameApi

import (
	"coinche/api"
	"coinche/domain"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinGame(test *testing.T) {
	assert := assert.New(test)
	mockGameService := MockGameService{
		map[int]domain.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO", Full: true},
		},
		nil,
	}
	router := api.SetupRouter(&mockGameService)

	test.Run("join game 1", func(test *testing.T) {
		request := testUtils.NewJoinGameRequest(1, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusAccepted, response.Code)
	})

	test.Run("should fail when joining non existing game", func(test *testing.T) {
		request := testUtils.NewJoinGameRequest(60, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusNotFound, response.Code)
	})

	test.Run("should fail when game is full", func(test *testing.T) {
		request := testUtils.NewJoinGameRequest(2, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusForbidden, response.Code)
	})
}
