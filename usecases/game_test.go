package usecases

import (
	"coinche/domain"
	testUtilities "coinche/utilities/test"
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
		gameID, err := gameUsecases.CreateGame("GAME TWO")
		testUtilities.FatalIfErr(err, test)

		assert.Equal(1, mockRepository.creationCalls)
		assert.Equal(1, gameID)
	})

	test.Run("can choose a team", func(test *testing.T) {
		err := gameUsecases.JoinTeam(0, "P1", "A Team")
		testUtilities.FatalIfErr(err, test)

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal("A Team", game.Players["P1"].Team)
	})
}
