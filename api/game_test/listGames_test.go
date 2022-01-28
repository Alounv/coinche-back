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

func TestListGames(test *testing.T) {
	assert := assert.New(test)
	gameService := MockGameService{
		map[int]domain.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO"},
		},
		nil,
	}
	router := api.SetupRouter(&gameService)

	test.Run("list games", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		response := httptest.NewRecorder()

		want := []domain.Game{
			{Name: "GAME ONE"},
			{Name: "GAME TWO"},
		}

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGames(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})
}
