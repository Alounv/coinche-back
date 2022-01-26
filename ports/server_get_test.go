package ports

import (
	"coinche/app"
	testUtils "coinche/utilities/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStore struct {
	names map[int]string
	setCalls []string
}

func (s *MockStore) GetAGame(id int) app.Game {
	game := app.Game{ Name: s.names[id], Id: id}
	return game
}

func (s *MockStore) CreateAGame(name string) int {
	s.setCalls = append(s.setCalls, name)
	return 1
}

func (s *MockStore) GetAllGames() []app.Game {
	var games []app.Game
	for id, name := range s.names {
		games = append(games, app.Game{ Name: name, Id: id})
	}
	return games
}

func TestServerGET(test *testing.T) {
	assert := assert.New(test)
	mockStore := MockStore{
		map[int]string{
		 	1: "GAME ONE",
			2: "GAME TWO",
		},
		nil,
	}
	router := SetupRouter(&mockStore)


	test.Run("get a game 1", func(test *testing.T) {
		want := app.Game(app.Game{Name:"GAME ONE", Id:1})

		request := testUtils.NewGETGameRequest(1)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})

	test.Run("get a game 2", func(t *testing.T) {
		want := app.Game(app.Game{Name:"GAME TWO", Id:2})

		request := testUtils.NewGETGameRequest(2)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		got := testUtils.DecodeToGame(response.Body, test)

		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(want, got)
	})

	test.Run("returns 404 on missing game", func(t *testing.T) {
		request := testUtils.NewGETGameRequest(3)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(http.StatusNotFound, response.Code)
	})
}

