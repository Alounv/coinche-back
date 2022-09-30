package api

import (
	"coinche/domain"
	"coinche/usecases"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletGame(test *testing.T) {
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
		request, _ := http.NewRequest(http.MethodDelete, "/games/1/delete", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusAccepted, response.Code)
	})
}
