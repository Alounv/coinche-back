package api

import (
	"coinche/domain"
	"coinche/usecases"
	"coinche/utilities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGame(test *testing.T) {
	assert := assert.New(test)
	mockRepository := usecases.NewMockGameRepo(
		map[int]*domain.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO", Players: map[string]domain.Player{
				"P1": {},
				"P2": {},
				"P3": {},
				"P4": {},
			}},
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)
	router, _ := SetupRouter(gameUsecases)

	test.Run("get a game 1", func(test *testing.T) {
		want := domain.Game(domain.Game{Name: "GAME ONE"})

		request := utilities.NewGetGameRequest(test, 1)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := utilities.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})

	test.Run("returns 404 on missing game", func(t *testing.T) {
		request := utilities.NewGetGameRequest(test, 3)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusNotFound, response.Code)
	})
}
