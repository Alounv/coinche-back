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
		map[int]domain.Game{1: game},
	)
	gameUsecases := NewGameUsecases(&mockRepository)

	test.Run("can join game", func(test *testing.T) {
		game, err := gameUsecases.JoinGame(1, "P4")

		assert.NoError(err)
		assert.Equal(4, len(game.Players))
		assert.Equal(domain.Teaming, game.Phase)
	})

	test.Run("can leave game", func(test *testing.T) {
		err := gameUsecases.LeaveGame(1, "P4")

		assert.NoError(err)

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(3, len(game.Players))
		assert.Equal(domain.Teaming, game.Phase)

		game, err = gameUsecases.JoinGame(1, "P4")

		assert.NoError(err)

		game, err = gameUsecases.GetGame(1)

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
		assert.Equal(2, gameID)

	})

	test.Run("can choose a team", func(test *testing.T) {
		err := gameUsecases.JoinTeam(1, "P1", "A Team")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal("A Team", game.Players["P1"].Team)
	})

	test.Run("can start a game", func(test *testing.T) {
		err := gameUsecases.JoinTeam(1, "P2", "A Team")
		if err != nil {
			test.Fatal(err)
		}
		err = gameUsecases.JoinTeam(1, "P3", "B Team")
		if err != nil {
			test.Fatal(err)
		}
		err = gameUsecases.JoinTeam(1, "P4", "B Team")
		if err != nil {
			test.Fatal(err)
		}

		err = gameUsecases.StartGame(1)
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(domain.Bidding, game.Phase)
	})

	test.Run("can bid", func(test *testing.T) {
		err := gameUsecases.Bid(1, "P1", 80, domain.Spade)
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

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
		err := gameUsecases.Pass(1, "P3")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(1, game.Bids[80].Pass)
		assert.Equal(0, game.Bids[80].Coinche)
	})

	test.Run("can coinche", func(test *testing.T) {
		err := gameUsecases.Coinche(1, "P1")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(0, game.Bids[80].Pass)
		assert.Equal(1, game.Bids[80].Coinche)
	})

	test.Run("can go to playing phase", func(test *testing.T) {
		err := gameUsecases.Pass(1, "P3")
		if err != nil {
			test.Fatal(err)
		}
		err = gameUsecases.Pass(1, "P4")
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(domain.Playing, game.Phase)
		assert.Equal(0, len(game.Deck))
		assert.Equal(8, len(game.Players["P1"].Hand))
		assert.Equal(8, len(game.Players["P2"].Hand))
		assert.Equal(8, len(game.Players["P3"].Hand))
		assert.Equal(8, len(game.Players["P4"].Hand))
	})

	test.Run("can play a card", func(test *testing.T) {
		playerHand := game.Players["P1"].Hand
		card := playerHand[0]
		err := gameUsecases.PlayCard(1, "P1", card)
		if err != nil {
			test.Fatal(err)
		}

		game, err := gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(7, len(game.Players["P1"].Hand))
		assert.Equal(8, len(game.Players["P2"].Hand))
		assert.Equal(8, len(game.Players["P3"].Hand))
		assert.Equal(8, len(game.Players["P4"].Hand))
	})

	test.Run("can archive a game", func(test *testing.T) {
		game, err := gameUsecases.GetGame(1)
		assert.Equal(1, game.Root)
		assert.NoError(err)

		err = gameUsecases.ArchiveGame(1)
		if err != nil {
			test.Fatal(err)
		}

		game, err = gameUsecases.GetGame(1)

		assert.NoError(err)
		assert.Equal(0, game.Root)
	})
}
