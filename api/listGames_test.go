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

func TestListGames(test *testing.T) {
	assert := assert.New(test)
	mockRepository := usecases.NewMockGameRepo(
		map[int]*domain.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO"},
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)
	router, _ := SetupRouter(gameUsecases)

	test.Run("list games", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := utilities.DecodeToGames(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Contains(got, domain.Game{ID: 1, Name: "GAME ONE"})
		assert.Contains(got, domain.Game{ID: 2, Name: "GAME TWO"})
	})
}
