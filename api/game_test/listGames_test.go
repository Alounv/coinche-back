package gameapi

import (
	"coinche/api"
	"coinche/domain"
	"coinche/usecases"
	testutils "coinche/utilities/test"
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
	router := api.SetupRouter(gameUsecases)

	test.Run("list games", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := testutils.DecodeToGames(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Contains(got, domain.Game{ID: 1, Name: "GAME ONE"})
		assert.Contains(got, domain.Game{ID: 2, Name: "GAME TWO"})
	})
}
