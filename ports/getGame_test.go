package ports

import (
	"coinche/app"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGame(test *testing.T) {
	assert := assert.New(test)
	mockStore := MockGameService{
		map[int]app.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO", Full: true},
		},
		nil,
	}
	router := SetupRouter(&mockStore)

	test.Run("get a game 1", func(test *testing.T) {
		want := app.Game(app.Game{Name: "GAME ONE"})

		request := testUtils.NewGetGameRequest(1)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})

	test.Run("get a game 2", func(t *testing.T) {
		want := app.Game(app.Game{Name: "GAME TWO", Full: true})

		request := testUtils.NewGetGameRequest(2)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})

	test.Run("returns 404 on missing game", func(t *testing.T) {
		request := testUtils.NewGetGameRequest(3)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusNotFound, response.Code)
	})
}
