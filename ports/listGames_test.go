package ports

import (
	"coinche/app"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListGames(test *testing.T) {
	assert := assert.New(test)
	mockStore := MockGameService{
		map[int]app.Game{
			1: {Name: "GAME ONE"},
			2: {Name: "GAME TWO"},
		},
		nil,
	}
	router := SetupRouter(&mockStore)

	test.Run("list games", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		response := httptest.NewRecorder()

		want := []app.Game{
			{Name: "GAME ONE"},
			{Name: "GAME TWO"},
		}

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGames(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})
}
