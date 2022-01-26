package ports

import (
	"coinche/app"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerGETList(test *testing.T) {
	assert := assert.New(test)
	mockStore := MockStore{
		map[int]string{
			1: "GAME ONE",
			2: "GAME TWO",
		},
		nil,
	}
	router := SetupRouter(&mockStore)

	test.Run("list games", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/games/all", nil)
		response := httptest.NewRecorder()

		want := []app.Game{
			{Id: 1, Name: "GAME ONE"},
			{Id: 2, Name: "GAME TWO"},
		}

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGames(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})
}
