package usecases

import (
	"coinche/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameService(test *testing.T) {
	assert := assert.New(test)

	mockRepository := MockGameRepo{
		games: map[int]*domain.Game{
			0: {Name: "GAME ONE", Players: map[string]domain.Player{
				"P1": {},
				"P2": {},
				"P3": {},
			}},
		},
	}
	gameUsecases := NewGameUsecases(&mockRepository)

	test.Run("can join game", func(test *testing.T) {
		game, err := gameUsecases.JoinGame(0, "P4")

		assert.NoError(err)
		assert.Equal(4, len(game.Players))
	})

	test.Run("can leave game", func(test *testing.T) {
		err := gameUsecases.LeaveGame(0, "P4")

		assert.NoError(err)

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(3, len(game.Players))
	})

	test.Run("can create game", func(test *testing.T) {
		id, err := gameUsecases.CreateGame("GAME TWO")
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(1, mockRepository.creationCalls)
		assert.Equal(1, id)
	})
}
