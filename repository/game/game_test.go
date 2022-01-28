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
		got := repository.GetGame(newID)

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

		got := repository.GetGame(1)

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
		want := []string{"P1", "P2", "P3", "P4"}

		err := repository.UpdateGame(2, want)
		if err != nil {
			panic(err)
		}
		got := repository.GetGame(2).Players

		assert.Equal(want, got)
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
