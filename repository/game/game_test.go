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

	repository := NewGameRepositoryFromDb(db)

	test.Run("create a game", func(test *testing.T) {
		newName := "NEW GAME ONE"
		newPlayers := []string{"P1", "P2"}

		newID := repository.CreateGame(domain.Game{Name: newName, Players: newPlayers})
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

	repository := NewGameRepositoryWithData(db)

	test.Run("get a game", func(test *testing.T) {
		want := domain.Game{Name: "GAME ONE", ID: 1, Players: []string{}}

		got, err := repository.GetGame(1)
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(want, got)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []domain.Game{
			{Name: "GAME ONE", ID: 1, Players: []string{}},
			{Name: "GAME TWO", ID: 2, Players: []string{"P1", "P2"}},
		}

		got := repository.ListGames()

		assert.Equal(want, got)
	})

	test.Run("update a game", func(test *testing.T) {
		players := []string{"P1", "P2", "P3", "P4"}

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

	test.Cleanup(func() {
		testutils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameRepositoryWithData(db *sqlx.DB) *GameRepository {
	repository := NewGameRepositoryFromDb(db)

	repository.CreateGames([]domain.Game{
		{Name: "GAME ONE", ID: 1},
		{Name: "GAME TWO", ID: 2, Players: []string{"P1", "P2"}},
	})

	return repository
}
