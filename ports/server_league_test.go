package ports

import (
	"coinche/app"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerGETLeague(test *testing.T) {
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

		router.ServeHTTP(response, request)

		want := []app.Game{
			{ Id: 1, Name: "GAME ONE"},
			{ Id: 2, Name: "GAME TWO"},
		}

		var got []app.Game
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Game, '%v'", response.Body, err)
		}
	
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})
}

