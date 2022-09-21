package usecases

import (
	"coinche/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameService(test *testing.T) {
	assert := assert.New(test)

	game := domain.NewGame("GAME ONE")
	game.Players = map[string]domain.Player{
		"P1": {},
		"P2": {},
		"P3": {},
	}
	mockRepository := NewMockGameRepo(
		map[int]domain.Game{0: game},
	)
	gameUsecases := NewGameUsecases(&mockRepository)

	test.Run("can join game", func(test *testing.T) {
		game, err := gameUsecases.JoinGame(0, "P4")

		assert.NoError(err)
		assert.Equal(4, len(game.Players))
		assert.Equal(domain.Teaming, game.Phase)
	})

	test.Run("can leave game", func(test *testing.T) {
		err := gameUsecases.LeaveGame(0, "P4")

		assert.NoError(err)

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(3, len(game.Players))
		assert.Equal(domain.Pause, game.Phase)

		game, err = gameUsecases.JoinGame(0, "P4")

		assert.NoError(err)

		game, err = gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(4, len(game.Players))
		assert.Equal(domain.Teaming, game.Phase)
	})

	test.Run("can create game", func(test *testing.T) {
		gameID, err := gameUsecases.CreateGame("GAME TWO")
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(1, mockRepository.creationCalls)
		assert.Equal(1, gameID)

	})

	test.Run("can choose a team", func(test *testing.T) {
		err := gameUsecases.JoinTeam(0, "P1", "A Team")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal("A Team", game.Players["P1"].Team)
	})

	test.Run("can start a game", func(test *testing.T) {
		err := gameUsecases.JoinTeam(0, "P2", "A Team")
		if err != nil {
			test.Fatal(err)
		}
		err = gameUsecases.JoinTeam(0, "P3", "B Team")
		if err != nil {
			test.Fatal(err)
		}
		err = gameUsecases.JoinTeam(0, "P4", "B Team")
		if err != nil {
			test.Fatal(err)
		}

		err = gameUsecases.StartGame(0)
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(domain.Bidding, game.Phase)
	})

	test.Run("can bid", func(test *testing.T) {
		err := gameUsecases.Bid(0, "P1", 80, domain.Spade)
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(domain.Spade, game.Bids[80].Color)
		assert.Equal("P1", game.Bids[80].Player)

		assert.Equal(0, game.Bids[80].Pass)
		assert.Equal(0, game.Bids[80].Coinche)

		assert.Equal(1, game.Players["P3"].Order)
		assert.Equal(2, game.Players["P2"].Order)
		assert.Equal(3, game.Players["P4"].Order)
		assert.Equal(4, game.Players["P1"].Order)
	})

	test.Run("can pass", func(test *testing.T) {
		err := gameUsecases.Pass(0, "P3")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(1, game.Bids[80].Pass)
		assert.Equal(0, game.Bids[80].Coinche)
	})

	test.Run("can coinche", func(test *testing.T) {
		err := gameUsecases.Coinche(0, "P1")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(0)

		assert.NoError(err)
		assert.Equal(0, game.Bids[80].Pass)
		assert.Equal(1, game.Bids[80].Coinche)
	})
}
