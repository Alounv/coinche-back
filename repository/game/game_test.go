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

	gameService := NewGameRepositaryFromDb(db)

	test.Run("create a game", func(test *testing.T) {
		newName := "NEW GAME ONE"

		newID := gameService.CreateGame(newName)
		got := gameService.GetGame(newID)

		assert.Equal(newName, got.Name)
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

	GameService := NewGameServiceWithData(db)

	test.Run("get a game", func(test *testing.T) {
		want := domain.Game{Name: "GAME ONE", ID: 1, Players: []string{}}

		got := GameService.GetGame(1)

		assert.Equal(want, got)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []domain.Game{
			{Name: "GAME ONE", ID: 1, Players: []string{}},
			{Name: "GAME TWO", ID: 2, Players: []string{"P1", "P2"}},
		}

		got := GameService.ListGames()

		assert.Equal(want, got)
	})

	test.Run("update a game", func(test *testing.T) {
		want := []string{"P1", "P2", "P3", "P4"}

		err := GameService.UpdateGame(2, want)
		if err != nil {
			panic(err)
		}
		got := GameService.GetGame(2).Players

		assert.Equal(want, got)
	})

	test.Cleanup(func() {
		testutils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameServiceWithData(db *sqlx.DB) *GameRepositary {
	dbGameService := NewGameRepositaryFromDb(db)

	dbGameService.CreateGames([]domain.Game{
		{Name: "GAME ONE", ID: 1},
		{Name: "GAME TWO", ID: 2, Players: []string{"P1", "P2"}},
	})

	return dbGameService
}
