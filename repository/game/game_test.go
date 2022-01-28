package gameRepo

import (
	"coinche/domain"
	"coinche/utilities/env"
	testUtils "coinche/utilities/test"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGameRepo(test *testing.T) {
	assert := assert.New(test)
	dbName := "testdb"
	env.LoadEnv("../../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtils.CreateDb(connectionInfo, dbName)

	gameService := NewGameRepoFromDb(db)

	test.Run("create a game", func(test *testing.T) {
		newName := "NEW GAME"

		newId := gameService.CreateGame(newName)
		got := gameService.GetGame(newId)

		assert.Equal(newName, got.Name)
		assert.Equal(newId, got.Id)
		assert.IsType(time.Time{}, got.CreatedAt)
	})

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}

func TestGameRepoWithInitialData(test *testing.T) {
	assert := assert.New(test)
	dbName := "testdb"
	env.LoadEnv("../../.env")
	connectionInfo := os.Getenv("SQLX_POSTGRES_INFO")

	db := testUtils.CreateDb(connectionInfo, dbName)

	GameService := NewGameServiceWithData(db)

	test.Run("get a game", func(test *testing.T) {
		want := domain.Game{Name: "GAME ONE", Id: 1, Players: []string{}}

		got := GameService.GetGame(1)

		assert.Equal(want, got)
	})

	test.Run("list all games", func(test *testing.T) {
		want := []domain.Game{
			{Name: "GAME ONE", Id: 1, Players: []string{}},
			{Name: "GAME TWO", Id: 2, Players: []string{"P1", "P2"}},
		}

		got := GameService.ListGames()

		assert.Equal(want, got)
	})

	test.Run("update a game", func(test *testing.T) {
		want := []string{"P1", "P2"}

		GameService.UpdateGame(2, want)
		got := GameService.GetGame(2).Players

		assert.Equal(want, got)
	})

	test.Cleanup(func() {
		testUtils.DropDb(connectionInfo, dbName, db)
	})
}

func NewGameServiceWithData(db *sqlx.DB) *GameRepo {
	dbGameService := NewGameRepoFromDb(db)

	dbGameService.CreateGames([]domain.Game{
		{Name: "GAME ONE", Id: 1},
		{Name: "GAME TWO", Id: 2, Players: []string{"P1", "P2"}},
	})

	return dbGameService
}
