package api

import (
	"coinche/domain"
	"coinche/usecases"
	testUtilities "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListGames(test *testing.T) {
	assert := assert.New(test)
	mockRepository := usecases.NewMockGameRepo(
		map[int]domain.Game{
			1: {
				Name:  "GAME ONE",
				Phase: domain.Teaming,
				Players: map[string]domain.Player{
					"P1": {},
					"P2": {},
					"P3": {},
				},
				Turns: []domain.Turn{{}, {}},
			},
			2: {
				Name:  "GAME TWO",
				Phase: domain.Teaming,
				Players: map[string]domain.Player{
					"P1": {},
					"P2": {},
					"P3": {},
				},
				Turns: []domain.Turn{{}, {}},
			},
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)
	router, _ := SetupRouter(gameUsecases)

	test.Run("list games", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := testUtilities.DecodeToGamePreviews(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(2, len(got))
		assert.Equal(3, len(got[1].Players))
		assert.Equal(domain.Teaming, got[1].Phase)
		assert.Equal(2, got[1].TurnsCount)
	})
}
