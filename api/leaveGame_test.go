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

func TestLeaveGame(test *testing.T) {
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
		},
	)
	gameUsecases := usecases.NewGameUsecases(&mockRepository)
	router, _ := SetupRouter(gameUsecases, []string{})

	test.Run("leave game", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/games/1/leave?playerName=P1", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusAccepted, response.Code)

		request = testUtilities.NewGetGameRequest(test, 1)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		got := testUtilities.DecodeToGame(response.Body, test)

		assert.Equal(2, len(got.Players))
	})
}
