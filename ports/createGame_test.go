package ports

import (
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGame(test *testing.T) {
	assert := assert.New(test)
	mockStore := MockGameService{
		map[int]string{},
		nil,
	}
	router := SetupRouter(&mockStore)

	test.Run("create a game", func(test *testing.T) {
		request := testUtils.NewPOSTGameRequest("NEW GAME")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusAccepted, response.Code)
		assert.Equal(1, len(mockStore.setCalls))
		assert.Equal("NEW GAME", mockStore.setCalls[0])
	})
}
