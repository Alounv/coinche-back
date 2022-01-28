package gameapi

import (
	"coinche/api"
	"coinche/domain"
	testutils "coinche/utilities/test"
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
			2: {Name: "GAME TWO", Players: []string{"P1", "P2", "P3", "P4"}},
		},
		nil,
	}
	router := api.SetupRouter(&mockGameService)

	test.Run("join game 1", func(test *testing.T) {
		request := testutils.NewJoinGameRequest(1, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusAccepted, response.Code)
	})

	test.Run("should fail when joining non existing game", func(test *testing.T) {
		request := testutils.NewJoinGameRequest(60, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusNotFound, response.Code)
	})

	test.Run("should fail when GAME IS FULL", func(test *testing.T) {
		request := testutils.NewJoinGameRequest(2, "Son Ly")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusForbidden, response.Code)
	})
}
