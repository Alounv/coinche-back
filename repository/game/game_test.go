package gamerepo

import (
	"coinche/domain"
	"coinche/utilities/env"
	testutils "coinche/utilities/test"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGameRepo(test *testing.T) {
	assert := assert.New(test)
	dbName := "testgamerepodb"
	env.LoadEnv("../../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testutils.CreateDb(connectionInfo, dbName)

	repository, err := NewGameRepositoryFromDb(db)
	if err != nil {
		test.Fatal(err)
	}

	test.Run("create a game", func(test *testing.T) {
		newName := "NEW GAME ONE"
		newPlayers := map[string]domain.Player{"P1": {}, "P2": {}}

		newID, err := repository.CreateGame(domain.Game{Name: newName, Players: newPlayers})
		if err != nil {
			test.Fatal(err)
		}

		got, err := repository.GetGame(newID)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(newName, got.Name)
		assert.Equal(newPlayers, got.Players)
		assert.Equal(newID, got.ID)
		assert.IsType(time.Time{}, got.CreatedAt)
	})

	test.Cleanup(func() {
		testutils.DropDb(connectionInfo, dbName, db)
	})
}

func TestGameRepoWithInitialData(test *testing.T) {
	assert := assert.New(test)
	dbName := "testgamerepowithinitialdatadb"
	env.LoadEnv("../../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testutils.CreateDb(connectionInfo, dbName)

	repository, err := NewGameRepositoryWithData(db)
	if err != nil {
		test.Fatal(err)
	}

	test.Run("get an empty game", func(test *testing.T) {
		want := domain.Game{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}}

		got, err := repository.GetGame(1)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want, got)
	})

	test.Run("get a game", func(test *testing.T) {
		want := domain.Game{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {}, "P2": {}}}

		got, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want, got)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []domain.Game{
			{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}},
			{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {}, "P2": {}}},
		}

		got, err := repository.ListGames()
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want[0], got[0])
		assert.Equal(want[1], got[1])
	})

	test.Run("update a game", func(test *testing.T) {
		players := map[string]domain.Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}}

		err := repository.UpdatePlayers(2, players, domain.Pause)
		if err != nil {
			panic(err)
		}
		game, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(players, game.Players)
		assert.Equal(domain.Pause, game.Phase)
	})

	test.Run("update a player", func(test *testing.T) {
		player := domain.Player{Team: "A Team"}

		err := repository.UpdatePlayer(2, "P2", player)
		if err != nil {
			panic(err)
		}
		game, err := repository.GetGame(2)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal("A Team", game.Players["P2"].Team)
	})

	test.Cleanup(func() {
		testutils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameRepositoryWithData(db *sqlx.DB) (*GameRepository, error) {
	repository, err := NewGameRepositoryFromDb(db)
	if err != nil {
		return nil, err
	}

	err = repository.CreateGames([]domain.Game{
		{Name: "GAME ONE", ID: 1, Players: map[string]domain.Player{}},
		{Name: "GAME TWO", ID: 2, Players: map[string]domain.Player{"P1": {}, "P2": {}}},
	})

	return repository, err
}
